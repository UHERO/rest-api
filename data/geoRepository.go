package data

import (
	"database/sql"

	"github.com/UHERO/rest-api/models"
)

type GeographyRepository struct {
	DB *sql.DB
}

func (r *FooRepository) GetGeographiesByCategory(categoryId int64) (geographies []models.DataPortalGeography, err error) {
	rows, err := r.DB.Query(
		`SELECT DISTINCT geographies.fips, geographies.handle, geographies.display_name, geographies.display_name_short, geographies.list_order
		FROM categories
		LEFT JOIN data_list_measurements ON data_list_measurements.data_list_id = categories.data_list_id
		LEFT JOIN measurement_series ON measurement_series.measurement_id = data_list_measurements.measurement_id
		LEFT JOIN series_v AS series ON series.id = measurement_series.series_id
		LEFT JOIN geographies ON geographies.id = series.geography_id
		WHERE (categories.id = ? OR categories.ancestry REGEXP CONCAT('[[:<:]]', ?, '[[:>:]]'))
		AND NOT (categories.hidden OR categories.masked)
		ORDER BY COALESCE(geographies.list_order, 999), geographies.handle`, categoryId, categoryId,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		var listOrder sql.NullInt64  // I hate that Go forces me to have this var to scan into, when SQL forces me to select it, but I don't need it :(
		geography := models.Geography{}
		err = rows.Scan(
			&geography.FIPS,
			&geography.Handle,
			&geography.Name,
			&geography.ShortName,
			&listOrder,
		)
		if err != nil {
			return
		}
		dataPortalGeography := models.DataPortalGeography{Handle: geography.Handle}
		if geography.FIPS.Valid {
			dataPortalGeography.FIPS = geography.FIPS.String
		}
		if geography.Name.Valid {
			dataPortalGeography.Name = geography.Name.String
		}
		if geography.ShortName.Valid {
			dataPortalGeography.ShortName = geography.ShortName.String
		}
		geographies = append(geographies, dataPortalGeography)
	}
	return
}

func (r *FooRepository) GetSeriesSiblingsGeoById(seriesId int64) (geographies []models.DataPortalGeography, err error) {
	rows, err := r.DB.Query(
		`SELECT DISTINCT geographies.fips, geographies.handle, geographies.display_name, geographies.display_name_short, geographies.list_order
		FROM series_v AS series
		JOIN (SELECT name, universe FROM series where id = ?) AS original_series  /* This "series" is base table, not confused with previous alias! */
		LEFT JOIN geographies ON geographies.id = series.geography_id
		WHERE series.universe = original_series.universe
		AND substring_index(series.name, '@', 1) = substring_index(original_series.name, '@', 1) /* prefixes are equal */
		ORDER BY COALESCE(geographies.list_order, 999), geographies.handle`, seriesId)
	if err != nil {
		return
	}
	for rows.Next() {
		var listOrder sql.NullInt64  // I hate that Go forces me to have this var to scan into, when SQL forces me to select it, but I don't need it :(
		geography := models.Geography{}
		err = rows.Scan(
			&geography.FIPS,
			&geography.Handle,
			&geography.Name,
			&geography.ShortName,
			&listOrder,
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
		if geography.ShortName.Valid {
			dataPortalGeography.ShortName = geography.ShortName.String
		}
		geographies = append(geographies, dataPortalGeography)
	}
	return
}
