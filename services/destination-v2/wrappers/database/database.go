package database

import (
	"context"
	"fmt"
	"github.com/bee-travels/bee-travels-go/destination-v2/wrappers/pgxpool"
	"github.com/elgris/sqrl"
	instana "github.com/instana/go-sensor"
	"github.com/pkg/errors"
	"os"
	"time"
)

func NewDatabasePool(sensor *instana.Sensor) (Pool, error) {
	host, found := os.LookupEnv("PG_HOST")
	if !found {
		return nil, errors.Errorf("PG_HOST must be set")
	}
	port, found := os.LookupEnv("PG_PORT")
	if !found {
		return nil, errors.Errorf("PG_PORT must be set")
	}
	user, found := os.LookupEnv("PG_USER")
	if !found {
		return nil, errors.Errorf("PG_USER must be set")
	}
	password, found := os.LookupEnv("PG_PASSWORD")
	if !found {
		return nil, errors.Errorf("PG_PASSWORD must be set")
	}
	connString := fmt.Sprintf(
		"host=%s port=%s dbname=beetravels user=%s password=%s",
		host, port, user, password,
	)

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
