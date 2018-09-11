package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetPackageRoutes(
	router *mux.Router,
	seriesRepository *data.SeriesRepository,
	searchRepository *data.SearchRepository,
	categoryRepository *data.CategoryRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router.HandleFunc(
		"/v1/package/series",
		controllers.GetSeriesPackage(seriesRepository, categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
		"cat", "{id:[0-9]+}",
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
		controllers.GetAnalyzerPackage(categoryRepository, seriesRepository, cacheRepository),
	).Methods("GET").Queries(
		"ids", "{ids_list:[0-9,]+}",
		"u", "{universe_text:.+}",
	)
	return router
}
