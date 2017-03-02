package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
)

func GetSeriesByCategoryId(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategory(id)
		returnSeriesList(seriesList, err, w, r)
	}
}

func GetInflatedSeriesByCategoryId(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetInflatedSeriesByCategory(id)
		returnInflatedSeriesList(seriesList, err, w, r)
	}
}

func GetSeriesByCategoryIdAndGeoHandle(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, ok := getIdAndGeo(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategoryAndGeo(id, geoHandle)
		returnSeriesList(seriesList, err, w, r)
	}
}

func GetSeriesByCategoryIdAndFreq(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, freq, ok := getIdAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategoryAndFreq(id, freq)
		returnSeriesList(seriesList, err, w, r)
	}
}

func GetSeriesByCategoryIdGeoHandleAndFreq(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategoryGeoAndFreq(id, geoHandle, freq)
		returnSeriesList(seriesList, err, w, r)
	}
}

func GetInflatedSeriesByCategoryIdGeoAndFreq(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetInflatedSeriesByCategoryGeoAndFreq(id, geoHandle, freq)
		returnInflatedSeriesList(seriesList, err, w, r)
	}
}

func GetSeriesById(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
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
		SetContext(r, j)
	}
}

func GetSeriesSiblingsById(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsById(id)
		returnSeriesList(seriesList, err, w, r)
	}
}

func GetSeriesSiblingsByIdAndFreq(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, freq, ok := getIdAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdAndFreq(id, freq)
		returnSeriesList(seriesList, err, w, r)
	}
}

func GetSeriesSiblingsByIdAndGeo(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geo, ok := getIdAndGeo(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdAndGeo(id, geo)
		returnSeriesList(seriesList, err, w, r)
	}
}

func GetSeriesSiblingsByIdGeoAndFreq(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geo, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdGeoAndFreq(id, geo, freq)
		returnSeriesList(seriesList, err, w, r)
	}
}

func GetSeriesSiblingsFreqById(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
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
		SetContext(r, j)
	}
}

func GetFreqByCategoryId(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
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
		SetContext(r, j)
	}
}

func GetSeriesObservations(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
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
		SetContext(r, j)
	}
}
