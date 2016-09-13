package routers

import (
	"github.com/gorilla/mux"
	"github.com/uhero/rest-api/data"
	"github.com/uhero/rest-api/controllers"
	"github.com/markbates/goth"
	"os"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/gothic"
	"github.com/codegangsta/negroni"
	"github.com/uhero/rest-api/common"
)

func SetApplicationRoutes(router *mux.Router, applicationRepository *data.ApplicationRepository) *mux.Router {
	goth.UseProviders(github.New(
		os.Getenv("GITHUB_KEY"),
		os.Getenv("GITHUB_SECRET"),
		"http://localhost:8080/auth/callback?provider=github",
	))
	router.HandleFunc("/auth/callback", controllers.AuthCallback(applicationRepository)).Methods("GET")
	router.HandleFunc("/auth", gothic.BeginAuthHandler).Methods("GET")
	router.HandleFunc("/", controllers.IndexHandler).Methods("GET")

	applicationRouter := mux.NewRouter()
	applicationRouter.HandleFunc("/applications", controllers.CreateApplication(applicationRepository)).Methods("POST")
	applicationRouter.HandleFunc("/applications", controllers.ReadApplications(applicationRepository)).Methods("GET")
	applicationRouter.HandleFunc("/applications/{id}", controllers.UpdateApplication(applicationRepository)).Methods("PUT", "POST")
	applicationRouter.HandleFunc("/applications/{id}", controllers.DeleteApplication(applicationRepository)).Methods("DELETE")
	router.PathPrefix("/applications").Handler(negroni.New(
		negroni.HandlerFunc(common.Authorize),
		negroni.Wrap(applicationRouter),
	))
	return router
}
