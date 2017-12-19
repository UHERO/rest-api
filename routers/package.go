package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetPackageRoutes(
	router *mux.Router,
	seriesRepository *data.SeriesRepository,
	categoryRepository *data.CategoryRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router.HandleFunc(
		"/v1/package/series",
		controllers.GetSeriesPackage(seriesRepository, categoryRepository, cacheRepository),
	).Methods("GET").Queries(
		"id", "{id:[0-9]+}",
	)
	router.HandleFunc(
		"/v1/package/search",
		controllers.GetSearchPackage(seriesRepository, cacheRepository),
	).Methods("GET").Queries(
		"q", "{search_text:.+}",
		"geo", "{geo:[A-Za-z-0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
		"u", "{universe_text:.+}",
	)
	return router
}
