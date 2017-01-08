package data

import (
	"database/sql"

	"github.com/UHERO/rest-api/models"
)

type GeographyRepository struct {
	DB *sql.DB
}

func (r *GeographyRepository) GetAllGeographies() (geographies []models.DataPortalGeography, err error) {
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
		dataPortalGeography := models.DataPortalGeography{Handle: geography.Handle}
		if geography.FIPS.Valid {
			dataPortalGeography.FIPS = geography.FIPS.String
		}
		if geography.Name.Valid {
			dataPortalGeography.Name = geography.Name.String
		}
		geographies = append(geographies, dataPortalGeography)
	}
	return
}

func (r *GeographyRepository) GetGeographiesByCategory(categoryId int64) (geographies []models.DataPortalGeography, err error) {
	rows, err := r.DB.Query(
		`SELECT geographies.fips, geographies.display_name_short, geographies.handle
  FROM
(SELECT DISTINCT(SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1)) AS handle
FROM categories
  LEFT JOIN data_list_measurements ON data_list_measurements.data_list_id = categories.data_list_id
  LEFT JOIN series ON series.measurement_id = data_list_measurements.measurement_id
WHERE categories.id = ? OR categories.ancestry REGEXP CONCAT('[[:<:]]', ?, '[[:>:]]') AND series.restricted = 0) AS category_handles
LEFT JOIN geographies ON geographies.handle LIKE category_handles.handle WHERE geographies.handle IS NOT NULL;`,
		categoryId,
		categoryId,
	)
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
		dataPortalGeography := models.DataPortalGeography{Handle: geography.Handle}
		if geography.FIPS.Valid {
			dataPortalGeography.FIPS = geography.FIPS.String
		}
		if geography.Name.Valid {
			dataPortalGeography.Name = geography.Name.String
		}
		geographies = append(geographies, dataPortalGeography)
	}
	return
}

func (r *GeographyRepository) GetSeriesSiblingsGeoById(seriesId int64) (geographies []models.DataPortalGeography, err error) {
	rows, err := r.DB.Query(`SELECT
		  geographies.fips, geographies.display_name_short,
		  catgeo.chandle AS handle
		FROM
		  (SELECT DISTINCT(SUBSTRING_INDEX(SUBSTR(catnames.name, LOCATE('@', catnames.name) + 1), '.', 1)) AS chandle
		    FROM
		      (SELECT series.name AS name
				FROM series JOIN (SELECT name FROM series where id = ?) as original_series
				WHERE series.name LIKE CONCAT(left(original_series.name, locate("@", original_series.name)), '%')
				AND series.restricted = 0)
				AS catnames) AS catgeo
			LEFT JOIN geographies ON catgeo.chandle LIKE geographies.handle;`, seriesId)
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
		dataPortalGeography := models.DataPortalGeography{Handle: geography.Handle}
		if geography.FIPS.Valid {
			dataPortalGeography.FIPS = geography.FIPS.String
		}
		if geography.Name.Valid {
			dataPortalGeography.Name = geography.Name.String
		}
		geographies = append(geographies, dataPortalGeography)
	}
	return
}
