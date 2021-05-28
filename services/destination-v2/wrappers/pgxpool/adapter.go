package pgxpool

import (
	"context"
	"errors"
	"fmt"
	instana "github.com/instana/go-sensor"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	ot "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net"
)

type Pool struct {
	sensor *instana.Sensor
	config *pgxpool.Config
	pool   *pgxpool.Pool
}

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Conn() *pgx.Conn
}

func Connect(sensor *instana.Sensor, ctx context.Context, connString string) (*Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	config.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		pid := conn.PgConn().PID()
		if err := conn.Ping(ctx); err != nil {
			fmt.Printf("Connection with PID %d seems bad, let's remove it from the pool\n", pid)
			return false
		}
		return true
	}

	return ConnectConfig(sensor, ctx, config)
}

func ConnectConfig(sensor *instana.Sensor, ctx context.Context, config *pgxpool.Config) (*Pool, error) {
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Pool{
		sensor: sensor,
		config: config,
		pool:   pool,
	}, nil
}

func ParseConfig(connString string) (*pgxpool.Config, error) {
	return pgxpool.ParseConfig(connString)
}

func (p *Pool) PgxPool() *pgxpool.Pool {
	return p.pool
}

func (p *Pool) Close() {
	p.pool.Close()
}

func (p *Pool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	childCtx, span := p.contextWithChildSpan(sql, ctx)
	defer span.Finish()

	tags, err := p.pool.Exec(childCtx, sql, args...)
	err = handleErr(err)
	if err != nil {
		span.SetTag(string(ext.Error), err.Error())
	}
	return tags, err
}

func (p *Pool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	childCtx, span := p.contextWithChildSpan(sql, ctx)
	defer span.Finish()

	rows, err := p.pool.Query(childCtx, sql, args...)
	err = handleErr(err)
	if err != nil {
		span.SetTag(string(ext.Error), err.Error())
	}
	return rows, err
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	childCtx, span := p.contextWithChildSpan(sql, ctx)
	conn, err := p.pool.Acquire(ctx)
	err = handleErr(err)
	return &singleRow{
		err:  err,
		conn: conn,
		sql:  sql,
		args: args,
		ctx:  childCtx,
		span: span,
	}
}

func (p *Pool) contextWithChildSpan(sql string, ctx context.Context) (context.Context, ot.Span) {
	var spanOptions []ot.StartSpanOption
	if parent, ok := instana.SpanFromContext(ctx); ok {
		spanOptions = append(spanOptions, ot.ChildOf(parent.Context()))
	}

	host := p.config.ConnConfig.Host
	port := p.config.ConnConfig.Port
	user := p.config.ConnConfig.User
	db := p.config.ConnConfig.Database

	span := p.sensor.Tracer().StartSpan(sql, spanOptions...)
	span.SetTag(string(ext.SpanKind), string(ext.SpanKindRPCClientEnum))
	span.SetTag(string(ext.DBType), "postgresql")
	span.SetTag(string(ext.DBInstance), db)
	span.SetTag(string(ext.DBUser), user)
	span.SetTag(string(ext.DBStatement), sql)
	span.SetTag(string(ext.PeerAddress), fmt.Sprintf("%s:%d", host, port))

	return instana.ContextWithSpan(ctx, span), span
}

func (p *Pool) Begin(ctx context.Context) (Tx, error) {
	t, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, handleErr(err)
	}
	return &tx{t: t, p: p}, nil
}

type tx struct {
	t pgx.Tx
	p *Pool
}

func (t *tx) Commit(ctx context.Context) error {
	childCtx, span := t.p.contextWithChildSpan("sql commit", ctx)
	defer span.Finish()

	err := t.t.Rollback(childCtx)
	err = handleErr(err)
	if err != nil {
		span.SetTag(string(ext.Error), err.Error())
	}
	return err
}

func (t *tx) Rollback(ctx context.Context) error {
	childCtx, span := t.p.contextWithChildSpan("sql rollback", ctx)
	defer span.Finish()

	err := t.t.Rollback(childCtx)
	err = handleErr(err)
	if err != nil {
		span.SetTag(string(ext.Error), err.Error())
	}
	return err
}

func (t *tx) Exec(ctx context.Context, sql string, args ...interface{}) (commandTag pgconn.CommandTag, err error) {
	childCtx, span := t.p.contextWithChildSpan(sql, ctx)
	defer span.Finish()

	tags, err := t.t.Exec(childCtx, sql, args...)
	err = handleErr(err)
	if err != nil {
		span.SetTag(string(ext.Error), err.Error())
	}
	return tags, err
}

func (t *tx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	childCtx, span := t.p.contextWithChildSpan(sql, ctx)
	defer span.Finish()

	rows, err := t.t.Query(childCtx, sql, args...)
	err = handleErr(err)
	if err != nil {
		span.SetTag(string(ext.Error), err.Error())
	}
	return rows, err
}

func (t *tx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	childCtx, span := t.p.contextWithChildSpan(sql, ctx)
	defer span.Finish()

	return rowWrapper{t.t.QueryRow(childCtx, sql, args...)}
}

func (t *tx) Conn() *pgx.Conn {
	return t.t.Conn()
}

type rowWrapper struct {
	row pgx.Row
}

func (r rowWrapper) Scan(dest ...interface{}) error {
	err := r.row.Scan(dest...)
	err = handleErr(err)
	return err
}

type singleRow struct {
	err  error
	conn *pgxpool.Conn
	sql  string
	args []interface{}
	ctx  context.Context
	span ot.Span
}

func (s *singleRow) Scan(dest ...interface{}) error {
	defer s.span.Finish()
	defer s.conn.Release()
	if s.err != nil {
		return handleErr(s.err)
	}
	row := s.conn.QueryRow(s.ctx, s.sql, s.args...)
	if err := row.Scan(dest...); err != nil {
		return handleErr(err)
	}
	return nil
}

func handleErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return context.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr != nil && netErr.Timeout() {
		return netErr
	}
	return err
}
