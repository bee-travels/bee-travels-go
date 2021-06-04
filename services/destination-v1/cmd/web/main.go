package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/bee-travels/bee-travels-go/services/destination-v1/internals/data"
	"github.com/bee-travels/bee-travels-go/services/destination-v1/internals/handler"
	"github.com/bee-travels/bee-travels-go/services/destination-v1/internals/service"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, time=${latency_human}\n",
	}))

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	d := loadData("./data/destinations.json")

	provider := service.NewLocalDB(d)

	h := handler.New(provider)

	v1 := e.Group("/api/v1")

	v1.GET("/destinations/:country/:city", h.GetDestinationByCityCountry)
	v1.GET("/destinations/:country", h.GetDestinationByCountry)
	v1.GET("/destinations", h.GetDestinations)

	e.Logger.Fatal(e.Start(":9001"))
}

func loadData(filename string) []data.Destination {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	var res []data.Destination

	err = json.Unmarshal(b, &res)

	if err != nil {
		log.Fatalln(err)
	}
	return res
}
