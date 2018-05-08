package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

// Proxy for census.gov
func SetCensusProxyRoute(
	router *mux.Router,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router.HandleFunc(
		`/v1/census/{census_endpoint:[a-zA-Z0-9=\-\/]+}`,
		controllers.GetCensusData(cacheRepository),
	).Methods("GET")
	return router
}
