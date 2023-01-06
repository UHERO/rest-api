package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetPackageRoutes(
	router *mux.Router,
	seriesRepository *data.FooRepository,
	searchRepository *data.SearchRepository,
	categoryRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router.HandleFunc(
		"/v1/package/series",
		controllers.GetSeriesPackage(seriesRepository, categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"cat", "{cat:[0-9]+}",
		"u", "{universe_text:.+}",
		"fc", "{forecast:[0-9Qq]+[FfHh](?:[0-9]+|[Ff])}",
	)
	router.HandleFunc(
		"/v1/package/series",
		controllers.GetSeriesPackage(seriesRepository, categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"cat", "{cat:[0-9]+}",
		"u", "{universe_text:.+}",
	)
	router.HandleFunc(
		"/v1/package/series",
		controllers.GetSeriesPackage(seriesRepository, categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"u", "{universe_text:.+}",
	)
	router.HandleFunc(
		"/v1/package/search",
		controllers.GetSearchPackage(searchRepository, cacheRepository),
	).Methods("GET").Queries(
		"q", "{search_text:.+}",
		"u", "{universe_text:.+}",
		"geo", "{geo:[A-Za-z-0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
	)
	router.HandleFunc(
		"/v1/package/search",
		controllers.GetSearchPackage(searchRepository, cacheRepository),
	).Methods("GET").Queries(
		"q", "{search_text:.+}",
		"u", "{universe_text:.+}",
	)
	router.HandleFunc(
		"/v1/package/category",
		controllers.GetCategoryPackage(categoryRepository, seriesRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"geo", "{geo:[A-Za-z0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
	)
	router.HandleFunc(
		"/v1/package/category",
		controllers.GetCategoryPackage(categoryRepository, seriesRepository, cacheRepository),
	).Methods("GET").Queries(
		"u", "{universe_text:.+}",
	)
	router.HandleFunc(
		"/v1/package/analyzer",
		controllers.GetAnalyzerPackage(categoryRepository, seriesRepository, cacheRepository, false),
	).Methods("GET").Queries(
		"ids", "{ids_list:[0-9,]+}",
		"u", "{universe_text:.+}",
	)
	router.HandleFunc(
		"/v1/package/analyzermom",
		controllers.GetAnalyzerPackage(categoryRepository, seriesRepository, cacheRepository, true),
	).Methods("GET").Queries(
		"ids", "{ids_list:[0-9,]+}",
		"u", "{universe_text:.+}",
	)
	router.HandleFunc(
		"/v1/package/export",
		controllers.GetExportPackage(seriesRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"expand", "{exp:[a-z]+}",
	)
	router.HandleFunc(
		"/v1/package/export",
		controllers.GetExportPackage(seriesRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)

	/* Following route exclusively for in-house staff use, returns unrestricted data. Only available to special in-house API instance */
	router.HandleFunc(
		"/v1.u/package/export",
		controllers.GetExportPackage(seriesRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"expand", "{exp:[a-z]+}",
	)
	router.HandleFunc(
		"/v1.u/package/export",
		controllers.GetExportPackage(seriesRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	return router
}
