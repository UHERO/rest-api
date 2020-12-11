package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetCategoryRoutes(
	router *mux.Router,
	categoryRepository *data.FooRepository,
	seriesRepository *data.FooRepository,
	measurementRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router.HandleFunc(
		"/v1/category",
		controllers.GetCategoryByIdGeoFreq(categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"geo", "{geo:[A-Za-z0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
	)
	router.HandleFunc(
		"/v1/category",
		controllers.GetCategory(categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	router.HandleFunc(
		"/v1/category",
		controllers.GetCategoriesByName(categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"search_text", "{searchText:.+}",
	)
	router.HandleFunc(
		"/v1/category",
		controllers.GetCategoriesByUniverse(categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"u", "{universe_text:.+}",
		"type", "{type_text:.+}",
	)
	router.HandleFunc(
		"/v1/category",
		controllers.GetCategoriesByUniverse(categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"u", "{universe_text:.+}",
	)
	router.HandleFunc(
		"/v1/category",
		controllers.GetCategoryRoots(categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"top_level", "true",
	)
	router.HandleFunc(
		"/v1/category",
		controllers.GetCategories(categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"top_level", "false",
	)
	router.HandleFunc(
		"/v1/category",
		controllers.GetCategories(categoryRepository, cacheRepository),
	).Methods("GET")

	router.HandleFunc(
		"/v1/category/freq",
		controllers.GetFreqByCategoryId(seriesRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	router.HandleFunc(
		"/v1/category/series",
		controllers.GetInflatedSeriesByGroupIdGeoAndFreq(seriesRepository, cacheRepository, data.Category),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"geo", "{geo:[A-Za-z0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
		"start", "{start_from:[12][0-9]{3}-[01][0-9]-[0-3][0-9]}",
		"expand", "true",
	)
	router.HandleFunc(
		"/v1/category/series",
		controllers.GetInflatedSeriesByGroupIdGeoAndFreq(seriesRepository, cacheRepository, data.Category),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"geo", "{geo:[A-Za-z0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
		"expand", "true",
	)
	router.HandleFunc(
		"/v1/category/series",
		controllers.GetSeriesByGroupIdGeoHandleAndFreq(seriesRepository, cacheRepository, data.Category),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"geo", "{geo:[A-Za-z0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
	)
	router.HandleFunc(
		"/v1/category/series",
		controllers.GetSeriesByGroupIdAndGeoHandle(seriesRepository, cacheRepository, data.Category),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"geo", "{geo:[A-Za-z0-9]+}",
	)
	router.HandleFunc(
		"/v1/category/series",
		controllers.GetSeriesByGroupIdAndFreq(seriesRepository, cacheRepository, data.Category),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
	)
	router.HandleFunc(
		"/v1/category/series",
		controllers.GetInflatedSeriesByGroupId(seriesRepository, cacheRepository, data.Category),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"expand", "true",
	)
	router.HandleFunc(
		"/v1/category/series",
		controllers.GetSeriesByGroupId(seriesRepository, cacheRepository, data.Category),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	router.HandleFunc(
		"/v1/category/measurements",
		controllers.GetMeasurementByCategoryId(measurementRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	return router
}
