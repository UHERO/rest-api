package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"errors"
)

func CreateFeedback(feedbackRepository *data.FeedbackRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var feedbackResource FeedbackResource
		err := json.NewDecoder(r.Body).Decode(&feedbackResource)
		if err != nil || &feedbackResource == nil {
			common.DisplayAppError(
				w,
				err,
				"Invalid feedback data",
				400,
			)
			return
		}
		feedback := feedbackResource.Data
		if feedback.Feedback == "" {
			common.DisplayAppError(
				w,
				errors.New("Empty feedback"),
				"Invalid feedback data",
				400,
			)
			return
		}
		_, err = feedbackRepository.CreateFeedback(feedback)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"Could not save feedback",
				500,
			)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}

