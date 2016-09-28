package routers

import (
	"github.com/gorilla/mux"
	"github.com/uhero/rest-api/data"
	"github.com/uhero/rest-api/controllers"
)

func SetCategoryRoutes(categoryRepository *data.CategoryRepository) *mux.Router {
	router := mux.NewRouter()
	categoryRouter := router.Path("/v1/category").Methods("GET").Subrouter()
	categoryRouter.Queries("id", "{id:[0-9]+").HandlerFunc(controllers.GetCategory(categoryRepository))
	categoryRouter.Queries("search_text", "{searchText:.+}").HandlerFunc(controllers.GetCategoriesByName(categoryRepository))
	categoryRouter.HandleFunc("", controllers.GetCategories(categoryRepository))
	return categoryRouter
}
