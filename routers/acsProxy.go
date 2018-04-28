package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/gorilla/mux"
)

// Proxy for ACS 5-Year Data (2009 - 2016)
func SetAcsProxyRoute(
	router *mux.Router,
) *mux.Router {
	router.HandleFunc(
		"/v1/acs",
		controllers.GetAcsData(),
	).Methods("GET").Queries(
		"get", "{ids_list:.+}",
	)
	return router
}
