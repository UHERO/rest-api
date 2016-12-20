package controllers

import (
	"encoding/json"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"net/http"
)

func CreateDataList(dataListRepository *data.DataListRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dataResource DataListResource
		// Decode the incoming application json
		err := json.NewDecoder(r.Body).Decode(&dataResource)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"Invalid Application Data",
				500,
			)
			return
		}
		//_, err = dataListRepository.CreateDataList()
	}
}
