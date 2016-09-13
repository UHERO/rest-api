package controllers

import (
	"github.com/uhero/rest-api/data"
	"net/http"
	"github.com/markbates/goth/gothic"
	"fmt"
	"html/template"
	"log"
	"github.com/uhero/rest-api/common"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"strconv"
)

func CreateApplication(creator data.Creator) func(http.ResponseWriter, *http.Request) {
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
		_, err = creator.Create(appClaims.Username, application)
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

func UpdateApplication(applicationRepository *data.ApplicationRepository) func(http.ResponseWriter, *http.Request) {
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

		_, err = applicationRepository.Update(appClaims.Username, application)
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

func ReadApplications(applicationRepository *data.ApplicationRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		appClaims, ok := common.FromContext(r.Context())
		if ok != true {
			panic(errors.New("cannot get value from context"))
		}
		applications, err := applicationRepository.GetAll(appClaims.Username)
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

func DeleteApplication(applicationRepository *data.ApplicationRepository) func(http.ResponseWriter, *http.Request) {
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
		_, err = applicationRepository.Delete(appClaims.Username, id)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Index Requested")
	t, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func AuthCallback(applicationRepository *data.ApplicationRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		applications, err := applicationRepository.GetAll(user.NickName)
		if err != nil {
			panic(err)
		}

		token, err := common.GenerateJWT(user.NickName, "user")
		if err != nil {
			panic(err)
		}

		userResource := UserResource{User: user, Applications: applications, Token: token}
		t, _ := template.New("userinfo").Parse(applicationTemplate)
		t.Execute(w, userResource)
	}
}

// View Templates
var indexTemplate = `
<p><a href="/auth?provider=github">Log in with GitHub</a></p>
`

var applicationTemplate = `
<h1>Applications for {{.User.NickName}} ({{.User.UserID}})</h1>
{{range .Applications}}
<h2>{{.Name}} ({{.Id}})</h2>
<dl class="dl-horizontal">
	<dt>Hostname</dt>
	<dd>{{.Hostname}}</dd>
	<dt>API Key</dt>
	<dd>{{.APIKey}}</dd>
</dl>
{{else}}
No Applications
{{end}}
<p>Token:</p>
<p>{{.Token}}</p>
<button>Add Application</button>
`
