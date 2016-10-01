package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/uhero/rest-api/common"
	"github.com/uhero/rest-api/data"
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
		if len(applications) == 0 {
			common.DisplayAppError(
				w,
				err,
				"Invalid API Key!",
				401,
			)
			return
		}
		log.Print("Off to the next handler")
		next(w, r)
	}
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
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
