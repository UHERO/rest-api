package routers

import (
	"github.com/gorilla/mux"
	"github.com/uhero/rest-api/controllers"
	"github.com/uhero/rest-api/data"
)

func SetCategoryRoutes(
	categoryRepository *data.CategoryRepository,
	seriesRepository *data.SeriesRepository,
) *mux.Router {
	categoryRouter := mux.NewRouter()
	categoryRouter.HandleFunc("/v1/category", controllers.GetCategory(categoryRepository)).Methods("GET").Queries("id", "{id:[0-9]+}")
	categoryRouter.HandleFunc("/v1/category", controllers.GetCategoriesByName(categoryRepository)).Methods("GET").Queries("search_text", "{searchText:.+}")
	categoryRouter.HandleFunc("/v1/category", controllers.GetCategoryRoots(categoryRepository)).Methods("GET").Queries("top_level", "true")
	categoryRouter.HandleFunc("/v1/category", controllers.GetCategories(categoryRepository)).Methods("GET").Queries("top_level", "false")
	categoryRouter.HandleFunc("/v1/category", controllers.GetCategories(categoryRepository)).Methods("GET")
	categoryRouter.HandleFunc("/v1/category/series", controllers.GetSeriesByCategoryId(seriesRepository)).Methods("GET").Queries("id", "{id:[0-9]+}")
	return categoryRouter
}
