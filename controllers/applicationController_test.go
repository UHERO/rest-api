package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/uhero/rest-api/common"
	"github.com/uhero/rest-api/data"
	"github.com/uhero/rest-api/models"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	applicationResult  = models.Application{Id: 1, Name: "foo", Hostname: "bar.com"}
	applicationsResult = []models.Application{
		applicationResult,
		{Id: 2, Name: "bar", Hostname: "foo.com"},
	}
)

type repoMock struct {
}

func (repoMock) Create(username string, application *models.Application) (int64, error) {
	return 1, nil
}

func (repoMock) Update(username string, application *models.Application) (int64, error) {
	return 1, nil
}

func (repoMock) Delete(username string, id int64) (int64, error) {
	return 1, nil
}

func (repoMock) GetAll(username string) ([]models.Application, error) {
	return applicationsResult, nil
}

type resultInterface interface {
}

type applicationTest struct {
	Function       func(data.Repository) func(http.ResponseWriter, *http.Request)
	RequestMethod  string
	RequestURL     string
	StatusCode     int
	ExpectedResult resultInterface
}

var applicationTests = []applicationTest{
	{
		CreateApplication,
		"POST",
		"/applications",
		http.StatusCreated,
		ApplicationResource{Data: applicationResult},
	},
	{
		ReadApplications,
		"GET",
		"/applications",
		http.StatusOK,
		ApplicationsResource{Data: applicationsResult},
	},
}

func TestApplicationController(t *testing.T) {
	mockResource := ApplicationResource{Data: applicationResult}

	j, err := json.Marshal(mockResource)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range applicationTests {
		individualTest(t, test, j)
	}
}

func individualTest(t *testing.T, test applicationTest, j []byte) {
	appClaims := &common.AppClaims{Username: "foobar"}
	req, err := http.NewRequest(test.RequestMethod, test.RequestURL, bytes.NewReader(j))
	if err != nil {
		t.Fatal(err)
	}
	reqCtx := req.WithContext(common.NewContext(context.Background(), appClaims))
	rr := httptest.NewRecorder()
	test.Function(repoMock{})(rr, reqCtx)

	// Check the status code is what we expect.
	if status := rr.Code; status != test.StatusCode {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, test.StatusCode)
	}
	// Check the value is what we expect.
	if reflect.TypeOf(test.ExpectedResult) == reflect.TypeOf(ApplicationResource{}) {
		singleApplicationCheck(t, test, rr)
		return
	}
	actualResult := ApplicationsResource{}
	err = json.NewDecoder(rr.Body).Decode(&actualResult)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(test.ExpectedResult, actualResult) {
		t.Fatal("Actual applications not equal to expected result")
	}
}

func singleApplicationCheck(t *testing.T, test applicationTest, rr *httptest.ResponseRecorder) {
	actualResult := ApplicationResource{}
	err := json.NewDecoder(rr.Body).Decode(&actualResult)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(test.ExpectedResult, actualResult) {
		t.Fatal("Actual application not equal to expected result")
	}
}
