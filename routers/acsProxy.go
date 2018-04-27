package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/gorilla/mux"
)

func SetAcsProxyRoute(
	router *mux.Router,
) *mux.Router {
	router.HandleFunc(
		"/v1/acs",
		controllers.GetAcsData(),
	).Methods("GET").Queries(
		"id", "{ids_list:.+}",
	)
	return router
}
