package main

import (
	"github.com/bee-travels/bee-travels-go/destination-v2/wrappers/database"
	"github.com/bee-travels/bee-travels-go/destination-v2/wrappers/server"
	"github.com/elgris/sqrl"
	"github.com/jackc/pgx/v4"
)

func queryLocations(pool database.Pool, ctx server.RequestContext, country string) ([]Location, error) {
	selector := buildBaseQuery(true)
	if country != "" {
		selector.Where("country = ?", country)
	}

	locations := make([]Location, 0)
	err := database.QueryFunc(pool, ctx, selector, func(row pgx.Row) error {
		var location Location
		err := row.Scan(&location.Country, &location.City)
		if err != nil {
			return err
		}

		locations = append(locations, location)
		return nil
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return locations, nil
		}
		return nil, err
	}
	return locations, nil
}

func queryDestinations(pool database.Pool, ctx server.RequestContext, country, city string) ([]Destination, error) {
	selector := buildBaseQuery(false).
		Where("country = ?", country).
		Where("city = ?", city)

	destinations := make([]Destination, 0)
	err := database.QueryFunc(pool, ctx, selector, func(row pgx.Row) error {
		var destination Destination
		err := row.Scan(
			&destination.ID,
			&destination.City,
			&destination.Country,
			&destination.Latitude,
			&destination.Longitude,
			&destination.Population,
			&destination.Description,
			&destination.Images,
		)
		if err != nil {
			return err
		}

		destinations = append(destinations, destination)
		return nil
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return destinations, nil
		}
		return nil, err
	}
	return destinations, nil
}

func buildBaseQuery(wildcardSelect bool) *sqrl.SelectBuilder {
	selector := database.QueryBuilder().Select()
	if wildcardSelect {
		selector.Columns("id", "city", "country", "latitude", "longitude", "population", "description", "images")
	} else {
		selector.Columns("country", "city")
	}
	return selector.From("destination")
}
