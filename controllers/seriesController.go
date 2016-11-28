package controllers

import (
	"encoding/json"
	"net/http"

	"errors"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
	"log"
)

func GetSeriesBySearchText(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
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
		seriesList, err := seriesRepository.GetSeriesBySearchText(searchText)
		returnSeriesList(seriesList, err, w)
	}
}

func GetSeriesByCategoryId(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategory(id)
		returnSeriesList(seriesList, err, w)
	}
}

func GetSeriesByCategoryIdAndGeoHandle(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, ok := getIdAndGeo(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategoryAndGeo(id, geoHandle)
		returnSeriesList(seriesList, err, w)
	}
}

func GetSeriesByCategoryIdAndFreq(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, freq, ok := getIdAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategoryAndFreq(id, freq)
		returnSeriesList(seriesList, err, w)
	}
}

func GetSeriesByCategoryIdGeoHandleAndFreq(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByCategoryGeoAndFreq(id, geoHandle, freq)
		returnSeriesList(seriesList, err, w)
	}
}

func GetSeriesById(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		log.Printf("Getting Series by id: %d", id)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)

	}
}

func GetSeriesSiblingsById(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		log.Printf("Getting Series by id: %d", id)
		seriesList, err := seriesRepository.GetSeriesSiblingsById(id)
		returnSeriesList(seriesList, err, w)
	}
}

func GetSeriesSiblingsByIdAndFreq(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, freq, ok := getIdAndFreq(w, r)
		if !ok {
			return
		}
		log.Printf("Getting Series Siblings by id and frequency: %d, %s", id, freq)
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdAndFreq(id, freq)
		returnSeriesList(seriesList, err, w)
	}
}

func GetSeriesSiblingsByIdAndGeo(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geo, ok := getIdAndGeo(w, r)
		if !ok {
			return
		}
		log.Printf("Getting Series Siblings by id and geo: %d, %s", id, geo)
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdAndGeo(id, geo)
		returnSeriesList(seriesList, err, w)
	}
}

func GetSeriesSiblingsByIdGeoAndFreq(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geo, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		log.Printf("Getting Series Siblings by id, geo, and freq: %d, %s, %s", id, geo, freq)
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdGeoAndFreq(id, geo, freq)
		returnSeriesList(seriesList, err, w)
	}
}

func GetSeriesSiblingsFreqById(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		log.Printf("Getting Series by id: %d", id)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}

func GetSeriesObservations(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		log.Printf("Getting Series by id: %d", id)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)

	}
}
