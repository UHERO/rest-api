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
	var hidden string
	if !r.ShowHiddenCats {
		hidden = `AND NOT (categories.hidden OR categories.masked)`
	}
	//language=MySQL
	query := `SELECT measurements.id, measurements.data_portal_name, data_list_measurements.indent
		FROM categories
		JOIN data_list_measurements ON categories.data_list_id = data_list_measurements.data_list_id
		JOIN measurements ON data_list_measurements.measurement_id = measurements.id
		WHERE categories.id = ?
		<%HIDDEN_COND%>
		ORDER BY data_list_measurements.list_order;`
	rows, err := r.RunQuery(ReplaceTemplateTag(query, "HIDDEN_COND", hidden), categoryId)
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
