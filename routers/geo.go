package routers

import (
	"github.com/gorilla/mux"
	"github.com/uhero/rest-api/controllers"
	"github.com/uhero/rest-api/data"
)

func SetGeographyRoutes(
	router *mux.Router,
	geoRepository *data.GeographyRepository,
) *mux.Router {
	router.HandleFunc("/v1/geo", controllers.GetGeographies(geoRepository)).Methods("GET")
	return router
}