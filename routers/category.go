package routers

import (
	"github.com/gorilla/mux"
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
)

func SetCategoryRoutes(
	router *mux.Router,
	categoryRepository *data.CategoryRepository,
	seriesRepository *data.SeriesRepository,
) *mux.Router {
	router.HandleFunc("/v1/category", controllers.GetCategory(categoryRepository)).Methods("GET").Queries("id", "{id:[0-9]+}")
	router.HandleFunc("/v1/category", controllers.GetCategoriesByName(categoryRepository)).Methods("GET").Queries("search_text", "{searchText:.+}")
	router.HandleFunc("/v1/category", controllers.GetCategoryRoots(categoryRepository)).Methods("GET").Queries("top_level", "true")
	router.HandleFunc("/v1/category", controllers.GetCategories(categoryRepository)).Methods("GET").Queries("top_level", "false")
	router.HandleFunc("/v1/category", controllers.GetCategories(categoryRepository)).Methods("GET")
	router.HandleFunc("/v1/category/series", controllers.GetSeriesByCategoryIdAndGeoHandle(seriesRepository)).Methods("GET").Queries("id", "{id:[0-9]+}", "geo", "{geo:[A-Za-z0-9]+}")
	router.HandleFunc("/v1/category/series", controllers.GetSeriesByCategoryId(seriesRepository)).Methods("GET").Queries("id", "{id:[0-9]+}")
	return router
}
