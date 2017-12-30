package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetSearchRoutes(
	router *mux.Router,
	searchRepository *data.SearchRepository,
	seriesRepository *data.SeriesRepository,
	cacheRepository *data.CacheRepository,
) *mux.Router {
	// deprecated
	router.HandleFunc("/v1/series", controllers.GetSeriesBySearchText(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"search_text", "{search_text:.+}",
	)
	// Ordering is important: More specific to less specific
	router.HandleFunc("/v1/search", controllers.GetSearchSummaryByUniverse(searchRepository, cacheRepository)).Methods("GET").Queries(
		"q", "{search_text:.+}",
		"u", "{universe_text:.+}",
	)
	router.HandleFunc("/v1/search", controllers.GetSearchSummary(searchRepository, cacheRepository)).Methods("GET").Queries(
		"q", "{search_text:.+}",
	)
	router.HandleFunc("/v1/search/series", controllers.GetInflatedSearchResultByGeoAndFreqAndUniverse(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"q", "{search_text:.+}",
		"geo", "{geo:[A-Za-z-0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
		"u", "{universe_text:.+}",
		"expand", "true",
	)
	router.HandleFunc("/v1/search/series", controllers.GetInflatedSearchResultByGeoAndFreq(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"q", "{search_text:.+}",
		"geo", "{geo:[A-Za-z-0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
		"expand", "true",
	)
	router.HandleFunc("/v1/search/series", controllers.GetSearchResultByGeoAndFreqAndUniverse(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"q", "{search_text:.+}",
		"geo", "{geo:[A-Za-z-0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
		"u", "{universe_text:.+}",
	)
	router.HandleFunc("/v1/search/series", controllers.GetSearchResultByGeoAndFreq(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"q", "{search_text:.+}",
		"geo", "{geo:[A-Za-z-0-9]+}",
		"freq", "{freq:[ASQMWDasqmwd]}",
	)
	router.HandleFunc("/v1/search/series", controllers.GetSeriesBySearchTextAndUniverse(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"q", "{search_text:.+}",
		"u", "{universe_text:.+}",
	)
	router.HandleFunc("/v1/search/series", controllers.GetSeriesBySearchText(seriesRepository, cacheRepository)).Methods("GET").Queries(
		"q", "{search_text:.+}",
	)
	return router
}
