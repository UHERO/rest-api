package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetGeographyRoutes(
	router *mux.Router,
	geoRepository *data.GeographyRepository,
) *mux.Router {
	router.HandleFunc("/v1/geo", controllers.GetGeographies(geoRepository)).Methods("GET")
	return router
}
