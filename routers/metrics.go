package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetMetricsRoutes(
	router *mux.Router,
	metricsRepository *data.MetricsRepository,
	applicationRepository *data.FooRepository,
) *mux.Router {
	// Current metrics endpoint
	router.HandleFunc("/v1/metrics", controllers.GetMetrics(metricsRepository)).Methods("GET")
	
	// Historical metrics endpoint with optional days parameter
	router.HandleFunc("/v1/metrics/historical", controllers.GetHistoricalMetrics(metricsRepository)).Methods("GET")
	
	return router
}