package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetSeriesRoutes(
	router *mux.Router,
	seriesRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router.HandleFunc("/v1/series", controllers.GetSeriesById(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	router.HandleFunc("/v1/series", controllers.GetSeriesByName(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"name", "{name:.+}",
		"u",	"{universe:[A-Za-z0-9]+}",
		"expand", "{exp:[a-z]+}",
	)
	router.HandleFunc("/v1/series", controllers.GetSeriesByName(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"name", "{name:.+}",
		"expand", "{exp:[a-z]+}",
	)
	router.HandleFunc("/v1/series", controllers.GetSeriesByName(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"name", "{name:.+}",
		"u",	"{universe:[A-Za-z0-9]+}",
	)
	router.HandleFunc("/v1/series", controllers.GetSeriesByName(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"name", "{name:.+}",
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

	/* Following routes exclusively for in-house staff use, return unrestricted data. Only available to special in-house API instance */
	router.HandleFunc("/v1.u/series", controllers.GetSeriesByName(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"name", "{name:.+}",
		"u",	"{universe:[A-Za-z0-9]+}",
		"start", "{start_from:[12][0-9]{3}-[01][0-9]-[0-3][0-9]}",
		"expand", "{exp:[a-z]+}",
	)
	router.HandleFunc("/v1.u/series", controllers.GetSeriesByName(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"name", "{name:.+}",
		"u",	"{universe:[A-Za-z0-9]+}",
		"expand", "{exp:[a-z]+}",
	)
	router.HandleFunc("/v1.u/series", controllers.GetSeriesByName(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"name", "{name:.+}",
		"start", "{start_from:[12][0-9]{3}-[01][0-9]-[0-3][0-9]}",
		"expand", "{exp:[a-z]+}",
	)
	router.HandleFunc("/v1.u/series", controllers.GetSeriesByName(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"name", "{name:.+}",
		"expand", "{exp:[a-z]+}",
	)
	router.HandleFunc("/v1.u/series", controllers.GetSeriesByName(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"name", "{name:.+}",
		"u",	"{universe:[A-Za-z0-9]+}",
	)
	router.HandleFunc("/v1.u/series", controllers.GetSeriesByName(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"name", "{name:.+}",
	)
	return router
}
