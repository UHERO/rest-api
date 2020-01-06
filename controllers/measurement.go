package controllers

import (
	"net/http"

	"encoding/json"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
)

func GetMeasurementByCategoryId(
	measurementRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		measurementList, err := measurementRepository.GetMeasurementsByCategory(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(MeasurementListResource{Data: measurementList})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error processing JSON has occurred",
				500,
			)
			return
		}
		WriteResponse(w, j)
		WriteCache(r, cacheRepository, j)
	}
}
