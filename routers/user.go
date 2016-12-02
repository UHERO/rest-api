package routers

import (
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/gplus"
	"os"
)

func SetUserRoutes(router *mux.Router, dataListRepository *data.DataListRepository) *mux.Router {
	goth.UseProviders(
		gplus.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), os.Getenv("GOOGLE_CALLBACK")),
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), os.Getenv("FACEBOOK_CALLBACK")),
	)
	router.HandleFunc("/user/auth/gplus/callback", controllers.ProviderCallback("gplus")).Methods("GET")
	router.HandleFunc("/user/auth/facebook/callback", controllers.ProviderCallback("facebook")).Methods("GET")
	router.HandleFunc("/user/auth", gothic.BeginAuthHandler).Methods("GET")

	userRouter := mux.NewRouter()
	userRouter.HandleFunc("/data_lists", controllers.CreateDataList(dataListRepository)).Methods("POST")
	userRouter.HandleFunc("/data_lists", controllers.ReadDataLists(dataListRepository)).Methods("GET")
	userRouter.HandleFunc("/data_lists/{id}", controllers.UpdateDataList(dataListRepository)).Methods("PUT", "POST")
	userRouter.HandleFunc("/data_lists/{id}", controllers.DeleteDataList(dataListRepository)).Methods("DELETE")
	router.PathPrefix("/data_lists/").Handler(negroni.New(
		negroni.HandlerFunc(common.IsAuthenticated),
		negroni.Wrap(userRouter),
	))
	return router
}
