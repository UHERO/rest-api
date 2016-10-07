package data

import (
	"database/sql"

	"github.com/UHERO/rest-api/models"
)

type GeographyRepository struct {
	DB *sql.DB
}

func (r *GeographyRepository) GetAllGeographies() (geographies []models.Geography, err error) {
	rows, err := r.DB.Query(`SELECT
	fips, display_name, handle FROM geographies;`)
	if err != nil {
		return
	}
	for rows.Next() {
		geography := models.Geography{}
		err = rows.Scan(
			&geography.FIPS,
			&geography.Name,
			&geography.Handle,
		)
		if err != nil {
			continue
		}
		geographies = append(geographies, geography)
	}
	return
}
