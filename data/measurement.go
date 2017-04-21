package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
)

type MeasurementRepository struct {
	DB *sql.DB
}

func (r *MeasurementRepository) GetMeasurementsByCategory(categoryId int64) (
	measurementList []models.Measurement,
	err error,
) {
	rows, err := r.DB.Query(`SELECT measurements.id, measurements.data_portal_name
		FROM categories
		LEFT JOIN data_list_measurements ON categories.data_list_id = data_list_measurements.data_list_id
		LEFT JOIN measurements ON data_list_measurements.measurement_id = measurements.id
		WHERE categories.id = ?
		AND NOT categories.hidden
		AND measurements.id IS NOT NULL;`,
		categoryId,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		measurement := models.Measurement{}
		err = rows.Scan(
			&measurement.Id,
			&measurement.Name,
		)
		if err != nil {
			return
		}
		measurementList = append(measurementList, measurement)
	}
	return
}
