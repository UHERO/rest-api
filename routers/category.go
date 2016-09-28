package routers

import (
	"github.com/gorilla/mux"
	"github.com/uhero/rest-api/data"
	"github.com/uhero/rest-api/controllers"
)

func SetCategoryRoutes(categoryRepository *data.CategoryRepository) *mux.Router {
	categoryRouter := mux.NewRouter()
	categoryRouter.HandleFunc("/v1/category", controllers.GetCategory(categoryRepository)).Methods("GET").Queries("id", "{id:[0-9]+}")
	categoryRouter.HandleFunc("/v1/category", controllers.GetCategoriesByName(categoryRepository)).Methods("GET").Queries("search_text", "{searchText:.+}")
	categoryRouter.HandleFunc("/v1/category", controllers.GetCategoryRoots(categoryRepository)).Methods("GET").Queries("top_level", "true")
	categoryRouter.HandleFunc("/v1/category", controllers.GetCategories(categoryRepository)).Methods("GET").Queries("top_level", "false")
	categoryRouter.HandleFunc("/v1/category", controllers.GetCategories(categoryRepository)).Methods("GET")
	return categoryRouter
}
