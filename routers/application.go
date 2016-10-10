package routers

import (
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"os"
)

func SetApplicationRoutes(router *mux.Router, applicationRepository *data.ApplicationRepository) *mux.Router {
	goth.UseProviders(github.New(
		os.Getenv("GITHUB_KEY"),
		os.Getenv("GITHUB_SECRET"),
		os.Getenv("GITHUB_CALLBACK"),
	))
	router.HandleFunc("/auth/callback", controllers.AuthCallback).Methods("GET")
	router.HandleFunc("/auth", gothic.BeginAuthHandler).Methods("GET")
	router.HandleFunc("/", controllers.IndexHandler).Methods("GET")

	applicationRouter := mux.NewRouter()
	applicationRouter.HandleFunc("/applications", controllers.CreateApplication(applicationRepository)).Methods("POST")
	applicationRouter.HandleFunc("/applications", controllers.ReadApplications(applicationRepository)).Methods("GET")
	applicationRouter.HandleFunc("/applications/{id}", controllers.UpdateApplication(applicationRepository)).Methods("PUT", "POST")
	applicationRouter.HandleFunc("/applications/{id}", controllers.DeleteApplication(applicationRepository)).Methods("DELETE")
	router.PathPrefix("/applications").Handler(negroni.New(
		negroni.HandlerFunc(common.IsAuthenticated),
		negroni.Wrap(applicationRouter),
	))
	//router.Handle("/developer", negroni.New(
	//	negroni.HandlerFunc(common.IsAuthenticated),
	//	negroni.Wrap(http.HandlerFunc(controllers.DeveloperHandler(applicationRepository))),
	//))
	return router
}
