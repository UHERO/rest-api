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
	router.HandleFunc("/v1/category/geo", controllers.GetGeographiesByCategory(geoRepository)).Methods("GET").Queries("id", "{id:[0-9]+}")
	return router
}
