package controllers

import (
	"encoding/json"
	"net/http"

	"errors"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
	"github.com/gorilla/context"
)

func GetSeriesBySearchText(searchRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := getSearchTextFromRequest(r)
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
		returnSeriesList(seriesList, err, w, r)
	}
}

func GetSearchSummary(searchRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := getSearchTextFromRequest(r)
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
		rUrl := r.URL.Path+"?"+r.URL.RawQuery
		context.Set(r, rUrl, j)
	}
}

func GetSearchResultByGeoAndFreq(searchRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := getSearchTextFromRequest(r)
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get searchText from request"),
				"Bad request.",
				400,
			)
			return
		}
		geo, ok := mux.Vars(r)["geo"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get geography from request"),
				"Bad request.",
				400,
			)
			return
		}
		freq, ok := mux.Vars(r)["freq"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get frequency from request"),
				"Bad request.",
				400,
			)
			return
		}
		seriesList, err := searchRepository.GetSearchResultsByGeoAndFreq(searchText, geo, freq)
		returnSeriesList(seriesList, err, w, r)
	}
}

func GetInflatedSearchResultByGeoAndFreq(searchRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := getSearchTextFromRequest(r)
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get searchText from request"),
				"Bad request.",
				400,
			)
			return
		}
		geo, ok := mux.Vars(r)["geo"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get geography from request"),
				"Bad request.",
				400,
			)
			return
		}
		freq, ok := mux.Vars(r)["freq"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get frequency from request"),
				"Bad request.",
				400,
			)
			return
		}
		seriesList, err := searchRepository.GetInflatedSearchResultsByGeoAndFreq(searchText, geo, freq)
		returnInflatedSeriesList(seriesList, err, w, r)
	}
}

func getSearchTextFromRequest(r *http.Request) (searchText string, ok bool) {
	searchText, ok = mux.Vars(r)["q"]
	if ok {
		return
	}
	searchText, ok = mux.Vars(r)["search_text"]
	return
}
