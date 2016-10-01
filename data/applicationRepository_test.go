package data

import (
	"github.com/uhero/rest-api/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
)

func TestCreateApplication(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO api_applications").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	username := "testUser"
	application := models.Application{
		Name:     "cool app",
		Hostname: "example.com",
	}

	applicationRepository := ApplicationRepository{DB: db}
	numRows, err := applicationRepository.CreateApplication(username, &application)

	if numRows != 1 {
		t.Fail()
	}
	if application.Id != 1 {
		t.Fail()
	}
	if len(application.APIKey) != 44 {
		t.Fail()
	}
}

func TestUpdateApplication(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("UPDATE api_applications").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	username := "testUser"
	application := models.Application{
		Id:       1,
		Name:     "cool app",
		Hostname: "example.com",
		APIKey:   "blah",
	}

	applicationRepository := ApplicationRepository{DB: db}
	numRows, err := applicationRepository.UpdateApplication(username, &application)
	if err != nil || numRows != 1 {
		t.Fail()
	}
}

func TestDeleteApplication(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("DELETE FROM api_applications").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	username := "testUser"
	application := models.Application{
		Id:       1,
		Name:     "cool app",
		Hostname: "example.com",
		APIKey:   "blah",
	}

	applicationRepository := ApplicationRepository{DB: db}
	numRows, err := applicationRepository.DeleteApplication(username, application.Id)
	if err != nil || numRows != 1 {
		t.Fail()
	}
}

func TestGetAllApplications(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	username := "testUser"
	application1 := models.Application{
		Id:       1,
		Name:     "cool app",
		Hostname: "example.com",
		APIKey:   "blah",
	}
	application2 := models.Application{
		Id:       2,
		Name:     "other cool app",
		Hostname: "example2.com",
		APIKey:   "blahblah",
	}
	applicationsResult := sqlmock.NewRows([]string{"id", "name", "hostname", "api_key"}).
		AddRow(application1.Id, application1.Name, application1.Hostname, application1.APIKey).
		AddRow(application2.Id, application2.Name, application2.Hostname, application2.APIKey)
	mock.ExpectQuery("SELECT id, name, hostname, api_key FROM api_applications WHERE github_nickname").
		WillReturnRows(applicationsResult)

	applicationRepository := ApplicationRepository{DB: db}
	applications, err := applicationRepository.GetAllApplications(username)
	if err != nil {
		t.Fail()
	}
	if len(applications) != 2 ||
		!reflect.DeepEqual(applications[0], application1) ||
		!reflect.DeepEqual(applications[1], application2) {
		t.Fail()
	}
}

func TestGetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	username := "testUser"
	application := models.Application{
		Id:       1,
		Name:     "cool app",
		Hostname: "example.com",
		APIKey:   "blah",
	}
	mockResult := sqlmock.NewRows([]string{"id", "name", "hostname", "api_key"}).
		AddRow(application.Id, application.Name, application.Hostname, application.APIKey)
	mock.ExpectQuery("SELECT id, name, hostname, api_key FROM api_applications").WillReturnRows(mockResult)

	applicationRepository := ApplicationRepository{DB: db}
	applicationResult, err := applicationRepository.GetApplicationById(username, application.Id)
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(applicationResult, application) {
		t.Fail()
	}
}
