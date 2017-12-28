package controllers

import (
	"encoding/json"
	"net/http"

	"errors"

	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func GetSearchPackage(
searchRepository *data.SeriesRepository,
cacheRepository *data.CacheRepository,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := mux.Vars(r)["search_text"]
		if !ok {
			common.DisplayAppError(w, errors.New("Couldn't get searchText from request"), "Bad request.", 400)
			return
		}
		universeText, ok := mux.Vars(r)["universe_text"]
		if !ok {
			common.DisplayAppError(w, errors.New("Couldn't get universeText from request"), "Bad request.", 400)
			return
		}
		var geo, freq string
		_, gotGeo := mux.Vars(r)["geo"]
		if gotGeo {
			geo, freq, ok = getGeoAndFreq(w, r)
			if !ok {
				return
			}
		}
		pkg, err := searchRepository.CreateSearchPackage(searchText, geo, freq, universeText)
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}

		j, err := json.Marshal(SearchPackage{Data: pkg})
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error processing JSON has occurred", 500)
			return
		}
		WriteResponse(w, j)
		WriteCache(r, cacheRepository, j)
	}
}

func GetSeriesBySearchText(searchRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := mux.Vars(r)["search_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get searchText from request"),
				"Bad request.",
				400,
			)
			return
		}
		seriesList, err := searchRepository.GetSeriesBySearchText(searchText)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetSeriesBySearchTextAndUniverse(searchRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := mux.Vars(r)["search_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get searchText from request"),
				"Bad request.",
				400,
			)
			return
		}
		universeText, ok := mux.Vars(r)["universe_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get universeText from request"),
				"Bad request.",
				400,
			)
			return
		}
		seriesList, err := searchRepository.GetSeriesBySearchTextAndUniverse(searchText, universeText)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetSearchSummary(searchRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := mux.Vars(r)["search_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get searchText from request"),
				"Bad request.",
				400,
			)
			return
		}
		searchSummary, err := searchRepository.GetSearchSummary(searchText)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(SearchSummaryResource{Data: searchSummary})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error processing JSON has occurred",
				500,
			)
			return
		}
		WriteResponse(w, j)
		WriteCache(r, c, j)
	}
}

func GetSearchSummaryByUniverse(searchRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := mux.Vars(r)["search_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get searchText from request"),
				"Bad request.",
				400,
			)
			return
		}
		universeText, ok := mux.Vars(r)["universe_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get universeText from request"),
				"Bad request.",
				400,
			)
			return
		}
		searchSummary, err := searchRepository.GetSearchSummaryByUniverse(searchText, universeText)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(SearchSummaryResource{Data: searchSummary})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error processing JSON has occurred",
				500,
			)
			return
		}
		WriteResponse(w, j)
		WriteCache(r, c, j)
	}
}

func GetSearchResultByGeoAndFreq(searchRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := mux.Vars(r)["search_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get searchText from request"),
				"Bad request.",
				400,
			)
			return
		}
		geo, freq, ok := getGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := searchRepository.GetSearchResultsByGeoAndFreq(searchText, geo, freq)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetSearchResultByGeoAndFreqAndUniverse(searchRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := mux.Vars(r)["search_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get searchText from request"),
				"Bad request.",
				400,
			)
			return
		}
		universeText, ok := mux.Vars(r)["universe_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get universeText from request"),
				"Bad request.",
				400,
			)
			return
		}
		geo, freq, ok := getGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := searchRepository.GetSearchResultsByGeoAndFreqAndUniverse(searchText, geo, freq, universeText)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetInflatedSearchResultByGeoAndFreq(searchRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := mux.Vars(r)["search_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get searchText from request"),
				"Bad request.",
				400,
			)
			return
		}
		geo, freq, ok := getGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := searchRepository.GetInflatedSearchResultsByGeoAndFreq(searchText, geo, freq)
		returnInflatedSeriesList(seriesList, err, w, r, c)
	}
}

func GetInflatedSearchResultByGeoAndFreqAndUniverse(searchRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := mux.Vars(r)["search_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get searchText from request"),
				"Bad request.",
				400,
			)
			return
		}
		universeText, ok := mux.Vars(r)["universe_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get universeText from request"),
				"Bad request.",
				400,
			)
			return
		}
		geo, freq, ok := getGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := searchRepository.GetInflatedSearchResultsByGeoAndFreqAndUniverse(searchText, geo, freq, universeText)
		returnInflatedSeriesList(seriesList, err, w, r, c)
	}
}
