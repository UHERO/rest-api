package controllers

import (
	"encoding/json"
	"net/http"

	"errors"
	"github.com/gorilla/mux"
	"github.com/uhero/rest-api/common"
	"github.com/uhero/rest-api/data"
	"strconv"
	"log"
)

func GetSeriesByCategoryId(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam, ok := mux.Vars(r)["id"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get category id from request"),
				"Bad request.",
				400,
			)
			return
		}
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}

		seriesList, err := seriesRepository.GetSeriesByCategory(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(SeriesListResource{Data: seriesList})
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

func GetSeriesById(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam, ok := mux.Vars(r)["id"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get series id from request"),
				"Bad request.",
				400,
			)
			return
		}
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
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

func GetSeriesObservations(seriesRepository *data.SeriesRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam, ok := mux.Vars(r)["id"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get series id from request"),
				"Bad request.",
				400,
			)
			return
		}
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
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
