package database

import (
	"context"
	"github.com/bee-travels/bee-travels-go/destination-v2/wrappers/pgxpool"
	"github.com/elgris/sqrl"
	instana "github.com/instana/go-sensor"
	"github.com/pkg/errors"
	"os"
	"time"
)

func NewDatabasePool(sensor *instana.Sensor) (Pool, error) {
	connString, found := os.LookupEnv("PG_CONN_STRING")
	if !found {
		return nil, errors.Errorf("PG_CONN_STRING must be set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	pool, err := pgxpool.Connect(sensor, ctx, connString)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func QueryBuilder() sqrl.StatementBuilderType {
	return sqrl.StatementBuilderType{}.PlaceholderFormat(sqrl.Dollar)
}
