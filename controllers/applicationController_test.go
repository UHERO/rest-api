package controllers

import (
	"testing"
	"github.com/uhero/rest-api/models"
	"net/http"
	"encoding/json"
	"net/http/httptest"
	"bytes"
	"context"
	"github.com/uhero/rest-api/common"
)

type repoMock struct {

}


func (repoMock) Create(username string, application *models.Application) (int64, error) {
	return 1, nil
}

func TestCreateApplication(t *testing.T) {
	createMock := repoMock{}
	mockApplication := models.Application{Id: 1, Name: "foo", Hostname: "bar.com"}
	mockResource := ApplicationResource{ Data: mockApplication}
	createHandler := CreateApplication(createMock)
	appClaims := &common.AppClaims{Username: "foobar"}

	j, err := json.Marshal(mockResource)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/applications", bytes.NewReader(j))
	if err != nil {
		t.Fatal(err)
	}
	reqCtx := req.WithContext(common.NewContext(context.Background(), appClaims))
	rr := httptest.NewRecorder()
	createHandler(rr, reqCtx)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	// Check the value is what we expect.
	var result ApplicationResource
	err = json.NewDecoder(rr.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Data.Name != mockApplication.Name {
		t.Fatalf("Expected Application Name %s, got %s", mockApplication.Name, result.Data.Name)
	}
}
