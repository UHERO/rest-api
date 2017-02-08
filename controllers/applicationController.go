package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
	"github.com/garyburd/redigo/redis"
	"strings"
)

const authPrefix = "Bearer "

// CreateApplication returns a handler that will create applications
func CreateApplication(applicationRepository data.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dataResource ApplicationResource
		// Decode the incoming Task json
		err := json.NewDecoder(r.Body).Decode(&dataResource)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"Invalid Application Data",
				500,
			)
			return
		}
		application := &dataResource.Data
		appClaims, ok := common.FromContext(r.Context())
		if ok != true {
			panic(errors.New("cannot get value from context"))
		}
		log.Printf("username: %s", appClaims.Username)
		_, err = applicationRepository.CreateApplication(appClaims.Username, application)
		if err != nil {
			panic(err)
		}
		j, err := json.Marshal(ApplicationResource{Data: *application})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(j)
	}
}

func ValidApiKey(applicationRepository *data.ApplicationRepository) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		authString := r.Header.Get("Authorization")
		if !strings.HasPrefix(authString, authPrefix) {
			common.DisplayAppError(
				w,
				errors.New("No Bearer Token"),
				"No Bearer Token!",
				401,
			)
			return
		}
		applications, err := applicationRepository.GetApplicationsByApiKey(strings.TrimPrefix(authString, authPrefix))
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"Error Veryfying API Key!",
				500,
			)
			return
		}
		if applications == nil || len(applications) == 0 {
			common.DisplayAppError(
				w,
				errors.New("Invalid API Key!"),
				"Invalid API Key!",
				401,
			)
			return
		}
		origin := r.Header.Get("Origin")
		if strings.HasPrefix(origin, "http://localhost") ||
			strings.Contains(origin, applications[0].Hostname) ||
			origin == "" {
			w.Header().Add("Access-Control-Allow-Origin", origin)
		}
		next(w, r)
	}
}

func CORSOptionsHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == http.MethodOptions {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Add("Access-Control-Allow-Headers", "authorization")
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
		return
	}
	next(w, r)
}

var RedisConn redis.Conn

func CheckCache() func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		url := r.URL.Path + "?" + r.URL.RawQuery
		cr, err := RedisConn.Do("GET", url)
		if err != nil {
			log.Fatal("Connection failure to Redis!")
			// Might want to rethink this
		}
		if cr == nil {
			log.Printf("Cache returned nil")
			next(w, r)
			return
		}
		log.Printf("Found |%s| in the cache", url)
		sendJSONResponseNoCache(w, cr.([]byte))
	}
}

func sendJSONResponseNoCache(w http.ResponseWriter, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func SendJSONResponse(w http.ResponseWriter, r *http.Request, payload []byte) {
	sendJSONResponseNoCache(w, payload)

	url := r.URL.Path + "?" + r.URL.RawQuery
	resp, err := RedisConn.Do("SET", url, payload)
	if err != nil {
		log.Fatal("Connection failure to Redis!")
		// Might want to rethink this
	}
	if resp != "OK" {
		log.Printf("DID NOT GET OK FROM REDIS")
	}
	log.Printf("Stored |%s| in the cache", url)
}

// UpdateApplication will return a handler for updating an application
func UpdateApplication(applicationRepository data.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dataResource ApplicationResource
		// Decode the incoming application json
		err := json.NewDecoder(r.Body).Decode(&dataResource)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"Invalid Application Data",
				500,
			)
			return
		}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			panic(err)
		}
		application := &dataResource.Data
		application.Id = id

		appClaims, ok := common.FromContext(r.Context())
		if ok != true {
			panic(errors.New("cannot get value from context"))
		}
		log.Printf("username: %s", appClaims.Username)

		_, err = applicationRepository.UpdateApplication(appClaims.Username, application)
		if err != nil {
			panic(err)
		}

		j, err := json.Marshal(ApplicationResource{Data: *application})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		SendJSONResponse(w, r, j)
	}
}

// ReadApplications returns a handler that returns all of the user's applications
func ReadApplications(applicationRepository data.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		appClaims, ok := common.FromContext(r.Context())
		if ok != true {
			panic(errors.New("cannot get value from context"))
		}
		applications, err := applicationRepository.GetAllApplications(appClaims.Username)
		if err != nil {
			panic(err)
		}
		j, err := json.Marshal(ApplicationsResource{Data: applications})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		SendJSONResponse(w, r, j)
	}
}

// DeleteApplication returns a handler that deletes an application
func DeleteApplication(applicationRepository data.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			panic(err)
		}

		appClaims, ok := common.FromContext(r.Context())
		if ok != true {
			panic(errors.New("cannot get value from context"))
		}
		_, err = applicationRepository.DeleteApplication(appClaims.Username, id)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
