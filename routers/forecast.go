package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetForecastRoutes(
	router *mux.Router,
	forecastRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router.HandleFunc("/v1/forecasts", controllers.GetForecasts(forecastRepository, cacheRepository)
	).Methods("GET").Queries(
		"which", "{which:(all|portal)}",
	)
	router.HandleFunc(
		"/v1/forecast/series",
		controllers.GetForecastSeries(forecastRepository, cacheRepository),
	).Methods("GET").Queries(
		"fc", "{forecast:[0-9Qq]+[FfHh](?:[0-9]+|[Ff])}",
		"freq", "{freq:[ASQMWDasqmwd]}",
	)
	router.HandleFunc(
		"/v1/forecast/series",
		controllers.GetForecastSeries(forecastRepository, cacheRepository),
	).Methods("GET").Queries(
		"fc", "{forecast:[0-9Qq]+[FfHh](?:[0-9]+|[Ff])}",
	)
	return router
}
