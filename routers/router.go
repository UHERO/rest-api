package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func InitRoutes(
	applicationRepository *data.FooRepository,
	categoryRepository *data.FooRepository,
	seriesRepository *data.FooRepository,
	searchRepository *data.SearchRepository,
	measurementRepository *data.FooRepository,
	geographyRepository *data.FooRepository,
	forecastRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	router = SetApplicationRoutes(router, applicationRepository)

	apiRouter := mux.NewRouter().StrictSlash(false)
	apiRouter = SetCategoryRoutes(apiRouter, categoryRepository, seriesRepository, measurementRepository, cacheRepository)
	apiRouter = SetSeriesRoutes(apiRouter, seriesRepository, cacheRepository)
	apiRouter = SetMeasurementRoutes(apiRouter, seriesRepository, cacheRepository)
	apiRouter = SetSearchRoutes(apiRouter, searchRepository, seriesRepository, cacheRepository)
	apiRouter = SetGeographyRoutes(apiRouter, geographyRepository, cacheRepository)
	apiRouter = SetPackageRoutes(apiRouter, seriesRepository, searchRepository, categoryRepository, cacheRepository)
	apiRouter = SetForecastRoutes(apiRouter, forecastRepository, cacheRepository)
	censusRouter := SetCensusProxyRoute(mux.NewRouter().StrictSlash(false), cacheRepository)

	router.PathPrefix("/v1/census").Handler(negroni.New(
		negroni.HandlerFunc(controllers.CORSOptionsHandler),
		negroni.HandlerFunc(controllers.ValidApiKey(applicationRepository)),
		negroni.HandlerFunc(controllers.CheckCacheFresh(cacheRepository)),
		negroni.Wrap(censusRouter),
	))
	router.PathPrefix("/v1").Handler(negroni.New(
		negroni.HandlerFunc(controllers.CORSOptionsHandler),
		negroni.HandlerFunc(controllers.ValidApiKey(applicationRepository)),
		negroni.HandlerFunc(controllers.CheckCache(cacheRepository)),
		negroni.Wrap(apiRouter),
	))
	return router
}
