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
	feedbackRepository *data.FeedbackRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	router = SetApplicationRoutes(router, applicationRepository)

	apiRouter := mux.NewRouter().StrictSlash(false)
	apiRouter = SetCategoryRoutes(apiRouter, categoryRepository, seriesRepository, cacheRepository)
	apiRouter = SetSeriesRoutes(apiRouter, seriesRepository, cacheRepository)
	apiRouter = SetSearchRoutes(apiRouter, seriesRepository, cacheRepository)
	apiRouter = SetGeographyRoutes(apiRouter, geographyRepository, cacheRepository)
	apiRouter = SetFeedbackRoutes(apiRouter, feedbackRepository)

	router.PathPrefix("/v1").Handler(negroni.New(
		negroni.HandlerFunc(controllers.CORSOptionsHandler),
		negroni.HandlerFunc(controllers.ValidApiKey(applicationRepository)),
		negroni.HandlerFunc(controllers.CheckCache(cacheRepository)),
		negroni.Wrap(apiRouter),
	))
	return router
}
