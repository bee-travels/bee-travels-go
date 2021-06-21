package main

import (
	"github.com/bee-travels/bee-travels-go/services/destination-v2/wrappers/database"
	"github.com/bee-travels/bee-travels-go/services/destination-v2/wrappers/server"
	instana "github.com/instana/go-sensor"
)

var lowercaseExceptions = []string{"es", "de", "au"}

func main() {
	if err := server.Start("destination-v2", initializeRouter); err != nil {
		panic(err)
	}
}

func initializeRouter(router server.PathRouter, pool database.Pool, _ *instana.Sensor) {
	router.Path("/api/v1/destinations", func(router server.PathRouter) {
		// path: /api/v1/destinations
		router.Get(listDestinations(pool))

		router.Path("/{country:string}", func(router server.PathRouter) {
			// path: /api/v1/destinations/:country
			router.Get(listDestinationsByCountry(pool))

			router.Path("/{city:string}", func(router server.PathRouter) {
				// path: /api/v1/destinations/:country/:city
				router.Get(listDestinationByCountryAndCity(pool))
			})
		})
	})
}
