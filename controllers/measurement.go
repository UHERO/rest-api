package controllers

import (
	"net/http"

	"github.com/UHERO/rest-api/data"
	"github.com/UHERO/rest-api/common"
	"encoding/json"
)

func GetMeasurementByCategoryId(measurementRepository *data.MeasurementRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
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
		SetCache(r, c, j)
	}
}

