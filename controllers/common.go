package controllers

import (
	"encoding/json"
	"errors"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/UHERO/rest-api/models"
	"github.com/gorilla/mux"
	"context"
	"net/http"
	"strconv"
	"log"
)

// Following for use by gorilla/context
type contextKey int
const cKey contextKey = 3

func CheckCache(c *data.CacheRepository) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		log.Printf("DEBUG: at entry CheckCache: r=%p", r)
		url := GetFullURL(r)
		cached_val, _ := c.GetCache(url)
		if cached_val == nil {
			log.Printf("DEBUG: Cache miss: url=%s", url)
			next(w, r)
			return
		}
		log.Printf("DEBUG: Cache HIT: "+url)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(cached_val)
	}
}

func SendJSONResponse(c *data.CacheRepository) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		log.Printf("DEBUG: at entry SendJSONResp: r is %p", r)
		url := GetFullURL(r)
		//if payload := FromContext(r.Context()); payload != nil { // }
		if payload, ok := context.GetOk(r, cKey); ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(payload.([]byte))
			err := c.SetCache(url, payload.([]byte))
			if err != nil {
				log.Printf("DEBUG: Cache store FAILURE: %s", url)
			} else {
				log.Printf("DEBUG: Stored in cache: %s", url)
			}
		} else {
			log.Printf("*** No data returned from context!")
		}
		next(w, r)
	}
}

func SetContext(r *http.Request, payload []byte) {
	context.Set(r, cKey, payload)
	return
}
/*
func NewContext(r *http.Request, payload []byte) context.Context {
	return context.WithValue(r.Context(), cKey, payload)
}

func FromContext(ctx context.Context) []byte {
	if ctx == nil {
		log.Printf("DEBUG: in FromCxt: ctx is nil")
		return nil
	} else {
		//log.Printf("DEBUG: in FromCxt: ctx NOT nil: %s", ctx)
	}
	//log.Printf("the payload is "+ctx.Value(cKey).(string))
	return ctx.Value(cKey).([]byte)
}
*/
func GetFullURL(r *http.Request) string {
	path := r.URL.Path
	if r.URL.RawQuery != "" {
		path += "?"+r.URL.RawQuery
	}
	return path
}

func returnSeriesList(seriesList []models.DataPortalSeries, err error, w http.ResponseWriter, r *http.Request) {
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
	if (string(j) == "foo") {
		log.Printf("bar")
	}
	//rUrl := r.URL.Path+"?"+r.URL.RawQuery
	//context.Set(r, rUrl, j)
}

func returnInflatedSeriesList(seriesList []models.InflatedSeries, err error, w http.ResponseWriter, r *http.Request) {
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
	if (string(j) == "foo") {
		log.Printf("bar")
	}
//	rUrl := r.URL.Path+"?"+r.URL.RawQuery
//	context.Set(r, rUrl, j)
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

