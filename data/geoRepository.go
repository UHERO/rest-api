package data

import (
	"database/sql"

	"github.com/UHERO/rest-api/models"
)

type GeographyRepository struct {
	DB *sql.DB
}

func (r *GeographyRepository) GetAllGeographies() (geographies []models.DataPortalGeography, err error) {
	rows, err := r.DB.Query(`SELECT fips, display_name, handle FROM geographies;`)
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
		`SELECT DISTINCT geographies.fips, geographies.display_name_short, geographies.handle
		FROM categories
		LEFT JOIN data_list_measurements ON data_list_measurements.data_list_id = categories.data_list_id
		LEFT JOIN measurement_series ON measurement_series.measurement_id = data_list_measurements.measurement_id
		LEFT JOIN series ON series.id = measurement_series.series_id
		     JOIN geographies ON geographies.id = series.geography_id
		LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
		WHERE (categories.id = ? OR categories.ancestry REGEXP CONCAT('[[:<:]]', ?, '[[:>:]]'))
		AND NOT categories.hidden
		AND NOT series.restricted
		AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined);`,
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
	rows, err := r.DB.Query(
		`SELECT DISTINCT geographies.fips, geographies.display_name_short, geographies.handle
		FROM series
		JOIN (SELECT name FROM series where id = ?) AS original_series
		JOIN geographies ON geographies.id = series.geography_id
		LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
		WHERE series.universe = 'UHERO'
		AND substring_index(series.name, '@', 1) = substring_index(original_series.name, '@', 1)
		AND NOT series.restricted
		AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined);`, seriesId)
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
