package main

import (
	"fmt"
	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
	"html/template"
	"log"
	"net/http"
	"os"
	"github.com/markbates/goth/providers/github"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uhero/rest-api/common"
)

func callbackAuthHandler(res http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}

	// here is where I have the user logged in and should have access
	// to the user.NickName user.UserID user.AccessToken
	// Step 1: get List of user applications from the applications repository

	//


	t, _ := template.New("userinfo").Parse(applicationTemplate)
	t.Execute(res, user)
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	log.Println("Index Requested")
	t, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	err = t.Execute(res, nil)
	if err != nil {
		fmt.Fprintln(res, err)
	}
}

func main() {
	// Set up MySQL
	connectionString := fmt.Sprintf(
		"%s:%s@%s(%s)/%s?parseTime=true&loc=US%%2FHawaii",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		common.AppConfig.Protocol,
		common.AppConfig.Server,
		common.AppConfig.Database,
	)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Register providers with Goth
	goth.UseProviders(
		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:8080/auth/twitter/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:8080/auth/github/callback"),
	)

	// Routing with Pat package
	r := pat.New()
	r.Get("/auth/{provider}/callback", callbackAuthHandler)
	r.Get("/auth/{provider}", gothic.BeginAuthHandler)
	r.Get("/", indexHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}

// View Templates
var indexTemplate = `
<p><a href="/auth/github">Log in with GitHub</a></p>
`

var applicationTemplate = `
<h1>Applications</h1>
{{range .Applications}}
<h2>{{.Name}}</h2>
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