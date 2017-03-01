package controllers

import (
	"encoding/json"
	"net/http"

	"errors"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
	"strconv"
)

func GetGeographies(geographyRepository *data.GeographyRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		geographies, err := geographyRepository.GetAllGeographies()
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(GeographiesResource{Data: geographies})
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
		//context.Set(r, rUrl, j)
	}
}

func GetGeographiesByCategory(geographyRepository *data.GeographyRepository) func(http.ResponseWriter, *http.Request) {
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
		geographies, err := geographyRepository.GetGeographiesByCategory(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(GeographiesResource{Data: geographies})
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
		//context.Set(r, rUrl, j)
	}
}

func GetSibllingGeographiesBySeriesId(geographyRepository *data.GeographyRepository) func(http.ResponseWriter, *http.Request) {
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
		geographies, err := geographyRepository.GetSeriesSiblingsGeoById(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(GeographiesResource{Data: geographies})
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
		//context.Set(r, rUrl, j)
	}
}
