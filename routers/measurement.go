package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetMeasurementRoutes(
	router *mux.Router,
	seriesRepository *data.SeriesRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router.HandleFunc(
		"/v1/measurement/series",
		controllers.GetInflatedSeriesByGroupIdGeoAndFreq(seriesRepository, cacheRepository, data.Measurement),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"geo", "{geo:[A-Za-z0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
		"expand", "true",
	)
	router.HandleFunc(
		"/v1/measurement/series",
		controllers.GetSeriesByGroupIdGeoHandleAndFreq(seriesRepository, cacheRepository, data.Measurement),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"geo", "{geo:[A-Za-z0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
	)
	router.HandleFunc(
		"/v1/measurement/series",
		controllers.GetSeriesByGroupIdAndGeoHandle(seriesRepository, cacheRepository, data.Measurement),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"geo", "{geo:[A-Za-z0-9]+}",
	)
	router.HandleFunc(
		"/v1/measurement/series",
		controllers.GetSeriesByGroupIdAndFreq(seriesRepository, cacheRepository, data.Measurement),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
	)
	router.HandleFunc(
		"/v1/measurement/series",
		controllers.GetInflatedSeriesByGroupId(seriesRepository, cacheRepository, data.Measurement),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"expand", "true",
	)
	router.HandleFunc(
		"/v1/measurement/series",
		controllers.GetSeriesByGroupId(seriesRepository, cacheRepository, data.Measurement),
	).Methods("GET").Queries("id", "{id:[0-9]+}")
	return router
}
