package database

import (
	"context"
	"fmt"
	"github.com/elgris/sqrl"
	instana "github.com/instana/go-sensor"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
	"github.com/kataras/iris/v12"
	ot "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"time"
)

type errorRow struct {
	e error
}

func (er errorRow) Scan(_ ...interface{}) error {
	return er.e
}

func Query(pool Pool, ctx iris.Context, query sqrl.Sqlizer) (pgx.Rows, error) {
	return QueryWithTimeout(pool, ctx, time.Second*20, query)
}

func QueryWithTimeout(pool Pool, ctx iris.Context, timeout time.Duration, query sqrl.Sqlizer) (pgx.Rows, error) {
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	stdCtx, span := contextWithChildSpan(sql, ctx.Request().Context())
	stdCtx, cancel := context.WithTimeout(stdCtx, timeout)

	rows, err := pool.Query(stdCtx, sql, args...)
	if err != nil {
		cancel()
		if err != pgx.ErrNoRows {
			span.SetTag(string(ext.Error), fmt.Sprintf("%+v", err))
			span.SetTag("params", args)
		}
		return nil, err
	}
	return pgRowsAdapter{span: span, rows: rows, cancel: cancel}, nil
}

func QueryRow(pool Pool, ctx iris.Context, query sqrl.Sqlizer) pgx.Row {
	return QueryRowWithTimeout(pool, ctx, time.Second*20, query)
}

func QueryRowWithTimeout(pool Pool, ctx iris.Context, timeout time.Duration, query sqrl.Sqlizer) pgx.Row {
	sql, args, err := query.ToSql()
	if err != nil {
		return errorRow{err}
	}

	stdCtx, span := contextWithChildSpan(sql, ctx.Request().Context())
	stdCtx, cancel := context.WithTimeout(stdCtx, timeout)

	row := pool.QueryRow(stdCtx, sql, args...)
	return pgRowAdapter{span: span, row: row, args: args, cancel: cancel}
}

func Exec(pool Pool, ctx iris.Context, query sqrl.Sqlizer) (pgconn.CommandTag, error) {
	return ExecWithTimeout(pool, ctx, time.Second*20, query)
}

func ExecWithTimeout(pool Pool, ctx iris.Context, timeout time.Duration, query sqrl.Sqlizer) (pgconn.CommandTag, error) {
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	stdCtx, span := contextWithChildSpan(sql, ctx.Request().Context())
	defer span.Finish()

	stdCtx, cancel := context.WithTimeout(stdCtx, timeout)
	defer cancel()

	return pool.Exec(stdCtx, sql, args...)
}

type RowFunction = func(row pgx.Row) error

func QueryFunc(pool Pool, ctx iris.Context, query sqrl.Sqlizer, fn RowFunction) error {
	return QueryFuncWithTimeout(pool, ctx, time.Second*20, query, fn)
}

func QueryFuncWithTimeout(pool Pool, ctx iris.Context, timeout time.Duration, query sqrl.Sqlizer, fn RowFunction) error {
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	stdCtx, span := contextWithChildSpan(sql, ctx.Request().Context())
	defer span.Finish()

	stdCtx, cancel := context.WithTimeout(stdCtx, timeout)
	defer cancel()

	rows, err := pool.Query(stdCtx, sql, args...)
	if err != nil {
		span.SetTag("params", args)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := fn(rows); err != nil {
			span.SetTag(string(ext.Error), fmt.Sprintf("%+v", err))
			span.SetTag("params", args)
			return err
		}
	}

	err = rows.Err()
	if err != nil {
		span.SetTag(string(ext.Error), fmt.Sprintf("%+v", err))
		span.SetTag("params", args)
	}
	return err
}

type pgRowAdapter struct {
	span   ot.Span
	row    pgx.Row
	cancel func()
	args   []interface{}
}

func (p pgRowAdapter) Scan(dest ...interface{}) error {
	err := p.row.Scan(dest...)
	if err != nil && err != pgx.ErrNoRows {
		p.span.SetTag(string(ext.Error), fmt.Sprintf("%+v", err))
		p.span.SetTag("params", p.args)
	}
	p.cancel()
	p.span.Finish()
	return err
}

type pgRowsAdapter struct {
	span   ot.Span
	rows   pgx.Rows
	cancel func()
}

func (p pgRowsAdapter) Close() {
	p.rows.Close()
	p.cancel()
	p.span.Finish()
}

func (p pgRowsAdapter) Err() error {
	return p.rows.Err()
}

func (p pgRowsAdapter) CommandTag() pgconn.CommandTag {
	return p.rows.CommandTag()
}

func (p pgRowsAdapter) FieldDescriptions() []pgproto3.FieldDescription {
	return p.rows.FieldDescriptions()
}

func (p pgRowsAdapter) Next() bool {
	return p.rows.Next()
}

func (p pgRowsAdapter) Scan(dest ...interface{}) error {
	return p.rows.Scan(dest...)
}

func (p pgRowsAdapter) Values() ([]interface{}, error) {
	return p.rows.Values()
}

func (p pgRowsAdapter) RawValues() [][]byte {
	return p.rows.RawValues()
}

func contextWithChildSpan(sql string, ctx context.Context) (context.Context, ot.Span) {
	var spanOptions []ot.StartSpanOption
	parent, ok := instana.SpanFromContext(ctx)
	if !ok {
		return ctx, nil
	}

	spanOptions = append(spanOptions, ot.ChildOf(parent.Context()))
	span := parent.Tracer().StartSpan(sql, spanOptions...)
	span.SetTag(string(ext.SpanKind), "intermediate")
	return instana.ContextWithSpan(ctx, span), span
}
