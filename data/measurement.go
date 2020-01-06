package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
)

type MeasurementRepository struct {
	DB *sql.DB
}

func (r *FooRepository) GetMeasurementsByCategory(categoryId int64) (
	measurementList []models.Measurement,
	err error,
) {
	//language=MySQL
	rows, err := r.RunQuery(`SELECT measurements.id, measurements.data_portal_name, data_list_measurements.indent
		FROM categories
		LEFT JOIN data_list_measurements ON categories.data_list_id = data_list_measurements.data_list_id
		LEFT JOIN measurements ON data_list_measurements.measurement_id = measurements.id
		WHERE categories.id = ?
		AND NOT (categories.hidden OR categories.masked)
		AND measurements.id IS NOT NULL
		ORDER BY data_list_measurements.list_order;`,
		categoryId,
	)

	if err != nil {
		return
	}
	for rows.Next() {
		measurement := models.Measurement{}
		indentString := sql.NullString{}
		err = rows.Scan(
			&measurement.Id,
			&measurement.Name,
			&indentString,
		)
		if err != nil {
			return
		}
		if indentString.Valid {
			measurement.Indent = indentationLevel[indentString.String]
		}
		measurementList = append(measurementList, measurement)
	}
	return
}
