package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/uhero/rest-api/controllers"
	"github.com/uhero/rest-api/data"
)

func InitRoutes(
	applicationRepository *data.ApplicationRepository,
	categoryRepository *data.CategoryRepository,
	seriesRepository *data.SeriesRepository,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	router = SetApplicationRoutes(router, applicationRepository)
	router.PathPrefix("/v1").Handler(negroni.New(
		negroni.HandlerFunc(controllers.ValidApiKey(applicationRepository)),
		negroni.Wrap(SetCategoryRoutes(categoryRepository, seriesRepository)),
	))
	return router
}
