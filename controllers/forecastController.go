package controllers

import (
	"encoding/json"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"net/http"
)

func GetForecasts(
	forecastRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pkg, err := forecastRepository.GetAllForecasts()
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
		forecast, ok := getStrParam(r, "forecast")
		if !ok {
			forecast = "@"  // a regex that will match any series name
		}
		seriesList, err := forecastRepository.GetForecastSeries(forecast)
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
		returnInflatedSeriesList(seriesList, err, w, r, cacheRepository)
	}
}