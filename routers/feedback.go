package routers

import (
	"github.com/UHERO/rest-api/controllers"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

func SetFeedbackRoutes(
	router *mux.Router,
	feedbackRepository *data.FeedbackRepository,
) *mux.Router {
	router.HandleFunc("/v1/feedback", controllers.CreateFeedback(feedbackRepository)).Methods("POST")
	return router
}
