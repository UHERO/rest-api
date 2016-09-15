package controllers

import (
	"fmt"
	"github.com/markbates/goth/gothic"
	"github.com/uhero/rest-api/common"
	"github.com/uhero/rest-api/data"
	"html/template"
	"log"
	"net/http"
	"errors"
)

// Landing page
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

// Authenticated Developer Landing Page
func DeveloperHandler(applicationRepository data.Repository) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {

		appClaims, ok := common.FromContext(r.Context())
		if ok != true {
			panic(errors.New("cannot get value from context"))
		}
		log.Printf("username: %s", appClaims.Username)
		applications, err := applicationRepository.GetAll(appClaims.Username)
		if err != nil {
			panic(err)
		}

		userResource := UserResource{User: appClaims.Username, Applications: applications}
		t, _ := template.New("userinfo").Parse(applicationTemplate)
		t.Execute(w, userResource)

	}
}

// Callback used by GitHub
func AuthCallback(w http.ResponseWriter, r *http.Request) {
	userProfile, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	token, err := common.GenerateJWT(userProfile.NickName, "user")
	common.StoreJWT(w, r, token)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/developer", http.StatusMovedPermanently)
}

// View Templates
var indexTemplate = `
<p><a href="/auth?provider=github">Log in with GitHub</a></p>
`

var applicationTemplate = `
<h1>Applications for {{.User}}</h1>
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
<button>Add Application</button>
`
