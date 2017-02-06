package controllers

import (
	"encoding/json"
	"errors"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func returnSeriesList(seriesList []models.DataPortalSeries, err error, w http.ResponseWriter) {
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
	SendJSONResponse(w, j)
}

func returnInflatedSeriesList(seriesList []models.InflatedSeries, err error, w http.ResponseWriter) {
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error has occurred",
			500,
		)
		return
	}
	j, err := json.Marshal(InflatedSeriesListResource{Data: seriesList})
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"An unexpected error processing JSON has occurred",
			500,
		)
		return
	}
	SendJSONResponse(w, j)
}

func getId(w http.ResponseWriter, r *http.Request) (id int64, ok bool) {
	ok = true
	idParam, gotId := mux.Vars(r)["id"]
	if !gotId {
		common.DisplayAppError(
			w,
			errors.New("Couldn't get id from request"),
			"Bad request.",
			400,
		)
		ok = false
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
		ok = false
		return
	}
	return
}

func getIdAndGeo(w http.ResponseWriter, r *http.Request) (id int64, geo string, ok bool) {
	ok = true
	idParam, gotId := mux.Vars(r)["id"]
	if !gotId {
		common.DisplayAppError(
			w,
			errors.New("Couldn't get category id from request"),
			"Bad request.",
			400,
		)
		ok = false
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
		ok = false
		return
	}
	geo, gotGeo := mux.Vars(r)["geo"]
	if !gotGeo {
		common.DisplayAppError(
			w,
			errors.New("Couldn't get geography handle from request"),
			"Bad request.",
			400,
		)
		ok = false
		return
	}
	return
}

func getIdAndFreq(w http.ResponseWriter, r *http.Request) (id int64, freq string, ok bool) {
	ok = true
	idParam, gotId := mux.Vars(r)["id"]
	if !gotId {
		common.DisplayAppError(
			w,
			errors.New("Couldn't get category id from request"),
			"Bad request.",
			400,
		)
		ok = false
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
		ok = false
		return
	}
	freq, gotFreq := mux.Vars(r)["freq"]
	if !gotFreq {
		common.DisplayAppError(
			w,
			errors.New("Couldn't get frequency from request"),
			"Bad request.",
			400,
		)
		ok = false
		return
	}
	return
}

func getIdGeoAndFreq(w http.ResponseWriter, r *http.Request) (id int64, geo string, freq string, ok bool) {
	ok = true
	idParam, gotId := mux.Vars(r)["id"]
	if !gotId {
		common.DisplayAppError(
			w,
			errors.New("Couldn't get category id from request"),
			"Bad request.",
			400,
		)
		ok = false
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
		ok = false
		return
	}
	geo, gotGeo := mux.Vars(r)["geo"]
	if !gotGeo {
		common.DisplayAppError(
			w,
			errors.New("Couldn't get geography handle from request"),
			"Bad request.",
			400,
		)
		ok = false
		return
	}
	freq, gotFreq := mux.Vars(r)["freq"]
	if !gotFreq {
		common.DisplayAppError(
			w,
			errors.New("Couldn't get frequency from request"),
			"Bad request.",
			400,
		)
		ok = false
		return
	}
	return
}
