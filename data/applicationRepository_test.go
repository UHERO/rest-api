package data

import (
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/uhero/rest-api/models"
)

func TestShouldCreateApplication(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO api_applications").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	userName := "testUser"
	application := models.Application{
		Name: "cool app",
		Hostname: "example.com",
	}

	applicationRepository := ApplicationRepository{DB: db}
	applicationRepository.Create(userName, &application)

	if application.Id != 1 {
		t.Fail()
	}
	if len(application.Key) != 44 {
		t.Fail()
	}
}
