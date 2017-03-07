package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
)

func GetSeriesByCategoryId(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategory(id)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetInflatedSeriesByCategoryId(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetInflatedSeriesByCategory(id)
		returnInflatedSeriesList(seriesList, err, w, r, c)
	}
}

func GetSeriesByCategoryIdAndGeoHandle(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, ok := getIdAndGeo(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategoryAndGeo(id, geoHandle)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetSeriesByCategoryIdAndFreq(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, freq, ok := getIdAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategoryAndFreq(id, freq)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetSeriesByCategoryIdGeoHandleAndFreq(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategoryGeoAndFreq(id, geoHandle, freq)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetInflatedSeriesByCategoryIdGeoAndFreq(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetInflatedSeriesByCategoryGeoAndFreq(id, geoHandle, freq)
		returnInflatedSeriesList(seriesList, err, w, r, c)
	}
}

func GetSeriesById(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		series, err := seriesRepository.GetSeriesById(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(SeriesResource{Data: series})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error processing JSON has occurred",
				500,
			)
			return
		}
		WriteResponseAndSetCache(w, r, c, j)
	}
}

func GetSeriesSiblingsById(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsById(id)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetSeriesSiblingsByIdAndFreq(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, freq, ok := getIdAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdAndFreq(id, freq)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetSeriesSiblingsByIdAndGeo(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geo, ok := getIdAndGeo(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdAndGeo(id, geo)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetSeriesSiblingsByIdGeoAndFreq(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geo, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdGeoAndFreq(id, geo, freq)
		returnSeriesList(seriesList, err, w, r, c)
	}
}

func GetSeriesSiblingsFreqById(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		frequencyList, err := seriesRepository.GetSeriesSiblingsFreqById(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(FrequencyListResource{Data: frequencyList})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error processing JSON has occurred",
				500,
			)
			return
		}
		WriteResponseAndSetCache(w, r, c, j)
	}
}

func GetFreqByCategoryId(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		frequencyList, err := seriesRepository.GetFreqByCategory(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(FrequencyListResource{Data: frequencyList})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error processing JSON has occurred",
				500,
			)
			return
		}
		WriteResponseAndSetCache(w, r, c, j)
	}
}

func GetSeriesObservations(seriesRepository *data.SeriesRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		series, err := seriesRepository.GetSeriesObservations(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(ObservationList{Data: series})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error processing JSON has occurred",
				500,
			)
			return
		}
		WriteResponseAndSetCache(w, r, c, j)
	}
}
