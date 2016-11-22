package routers

import (
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/facebook"
	"os"
	"github.com/UHERO/rest-api/controllers"
)

func SetUserRoutes(router *mux.Router) *mux.Router {
	goth.UseProviders(
		gplus.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), os.Getenv("GOOGLE_CALLBACK")),
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), os.Getenv("FACEBOOK_CALLBACK")),
	)
	router.HandleFunc("/user/auth/gplus/callback", controllers.ProviderCallback("gplus")).Methods("GET")
	router.HandleFunc("/user/auth/facebook/callback", controllers.ProviderCallback("facebook")).Methods("GET")
	router.HandleFunc("/user/auth", gothic.BeginAuthHandler).Methods("GET")
	return router
}