package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetGeographyRoutes(
	router *mux.Router,
	geoRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router.HandleFunc("/v1/category/geo", controllers.GetGeographiesByCategory(geoRepository, cacheRepository)).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"fc", "{forecast:[0-9Qq]+[FfHh](?:[0-9]+|[Ff])}",
	)
	router.HandleFunc("/v1/category/geo", controllers.GetGeographiesByCategory(geoRepository, cacheRepository)).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	router.HandleFunc("/v1/series/siblings/geo", controllers.GetSibllingGeographiesBySeriesId(geoRepository, cacheRepository)).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	return router
}
