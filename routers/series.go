package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetSeriesRoutes(
	router *mux.Router,
	seriesRepository *data.SeriesRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router.HandleFunc("/v1/series", controllers.GetSeriesById(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	router.HandleFunc("/v1/series", controllers.GetSeriesByName(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"name", "{name:[0-9]+}",
	)
	router.HandleFunc("/v1/series/siblings", controllers.GetSeriesSiblingsByIdGeoAndFreq(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"geo", "{geo:[A-Za-z0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
	)
	router.HandleFunc("/v1/series/siblings", controllers.GetSeriesSiblingsByIdAndGeo(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"geo", "{geo:[A-Za-z0-9]+}",
	)
	router.HandleFunc("/v1/series/siblings", controllers.GetSeriesSiblingsByIdAndFreq(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
	)
	router.HandleFunc("/v1/series/siblings", controllers.GetSeriesSiblingsById(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	router.HandleFunc("/v1/series/siblings/freq", controllers.GetSeriesSiblingsFreqById(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	router.HandleFunc("/v1/series/observations", controllers.GetSeriesObservations(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	return router
}
