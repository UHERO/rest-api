package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func InitRoutes(
	applicationRepository *data.ApplicationRepository,
	categoryRepository *data.CategoryRepository,
	seriesRepository *data.SeriesRepository,
	geographyRepository *data.GeographyRepository,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	router = SetApplicationRoutes(router, applicationRepository)
	router = SetUserRoutes(router)

	apiRouter := mux.NewRouter().StrictSlash(false)
	apiRouter = SetCategoryRoutes(apiRouter, categoryRepository, seriesRepository)
	apiRouter = SetSeriesRoutes(apiRouter, seriesRepository)
	apiRouter = SetGeographyRoutes(apiRouter, geographyRepository)

	router.PathPrefix("/v1").Handler(negroni.New(
		negroni.HandlerFunc(controllers.CORSOptionsHandler),
		negroni.HandlerFunc(controllers.ValidApiKey(applicationRepository)),
		negroni.Wrap(apiRouter),
	))
	return router
}
