package routers

import (
	"github.com/gorilla/mux"
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
)

func SetSeriesRoutes(
	router *mux.Router,
	seriesRepository *data.SeriesRepository,
) *mux.Router {
	router.HandleFunc("/v1/series", controllers.GetSeriesById(seriesRepository)).Methods("GET").Queries("id", "{id:[0-9]+}")
	router.HandleFunc("/v1/series/observations", controllers.GetSeriesObservations(seriesRepository)).Methods("GET").Queries("id", "{id:[0-9]+}")
	return router
}
