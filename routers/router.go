package routers

import (
	"github.com/gorilla/mux"
	"github.com/uhero/rest-api/data"
)

func InitRoutes(applicationRepository *data.ApplicationRepository) *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	// Routes for the User entity
	router = SetApplicationRoutes(router, applicationRepository)
	return router
}
