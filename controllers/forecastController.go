package controllers

import (
	"encoding/json"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/UHERO/rest-api/models"
	"net/http"
)

func GetForecasts(
	forecastRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		which, ok := getStrParam(r, "which")
		if !ok {
			which = "all"
		}
		var pkg models.ForecastList
		var err error
		if which == "portal" {
			pkg, err = forecastRepository.GetAllPortalForecasts()
		} else {
			pkg, err = forecastRepository.GetAllForecasts()
		}
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
		var j []byte
		j, err = json.Marshal(ForecastsResource{Data: pkg})
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error processing JSON has occurred", 500)
			return
		}
		WriteResponse(w, j)
		WriteCache(r, cacheRepository, j)
	}
}

func GetForecastSeries(
	forecastRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		freq, ok := getStrParam(r, "freq")
		if !ok {
			freq = ""
		}
		forecast, ok := getStrParam(r, "forecast")
		if !ok {
			forecast = "@"  // a regex that will match any series name
		}
		seriesList, err := forecastRepository.GetForecastSeries(forecast, freq)
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
		returnInflatedSeriesList(seriesList, err, w, r, cacheRepository)
	}
}
