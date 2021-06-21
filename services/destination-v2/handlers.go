package main

import (
	"github.com/bee-travels/bee-travels-go/services/destination-v2/wrappers/database"
	"github.com/bee-travels/bee-travels-go/services/destination-v2/wrappers/server"
	"net/http"
)

func listDestinations(pool database.Pool) server.RequestHandler {
	return func(ctx server.RequestContext) {
		location, err := queryLocations(pool, ctx, "")
		if err != nil {
			server.Response(ctx, http.StatusForbidden, Error{
				Error: err.Error(),
			})
			return
		}
		server.Response(ctx, http.StatusOK, location)
	}
}

func listDestinationsByCountry(pool database.Pool) server.RequestHandler {
	return func(ctx server.RequestContext) {
		country := ctx.Params().Get("country")
		location, err := queryLocations(pool, ctx, capitalize(country))
		if err != nil {
			server.Response(ctx, http.StatusForbidden, Error{
				Error: err.Error(),
			})
			return
		}
		server.Response(ctx, http.StatusOK, location)
	}
}

func listDestinationByCountryAndCity(pool database.Pool) server.RequestHandler {
	return func(ctx server.RequestContext) {
		country := ctx.Params().Get("country")
		city := ctx.Params().Get("city")
		destinations, err := queryDestinations(pool, ctx, capitalize(country), capitalize(city))
		if err != nil {
			server.Response(ctx, http.StatusForbidden, Error{
				Error: err.Error(),
			})
			return
		}
		server.Response(ctx, http.StatusOK, destinations)
	}
}
