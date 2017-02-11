package data

import (
	"github.com/UHERO/rest-api/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestCreateFeedback(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO user_feedbacks").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	feedback := models.Feedback{
		Name: "Peter Piper",
		Email: "picked@pickle.com",
		Feedback: "Awesome website",
	}

	feedbackRepository := FeedbackRepository{DB: db}
	numRows, err := feedbackRepository.CreateFeedback(feedback)
	if err != nil || numRows != 1 {
		t.Fail()
	}
}
