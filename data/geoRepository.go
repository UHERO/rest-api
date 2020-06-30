package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
)

type GeographyRepository struct {
	DB *sql.DB
}

func (r *FooRepository) GetGeographiesByCategory(categoryId int64) (geographies []models.DataPortalGeography, err error) {
	//language=MySQL
	rows, err := r.RunQuery(
		`SELECT DISTINCT geo_fips, geo_handle, geo_display_name, geo_display_name_short
		FROM <%PORTAL%> pv
		WHERE (category_id = ? OR category_ancestry REGEXP CONCAT('[[:<:]]', ?, '[[:>:]]'))
		ORDER BY COALESCE(geo_list_order, 999), geo_handle`, categoryId, categoryId,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		geography := models.Geography{}
		err = rows.Scan(
			&geography.FIPS,
			&geography.Handle,
			&geography.Name,
			&geography.ShortName,
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
	//language=MySQL
	rows, err := r.RunQuery(
		`SELECT DISTINCT geographies.fips, geographies.handle, geographies.display_name, geographies.display_name_short
		FROM <%SERIES%> AS series
		JOIN (SELECT name, universe FROM <%SERIES%> where id = ?) AS original_series
		JOIN geographies ON geographies.id = series.geography_id
		WHERE series.universe = original_series.universe
		AND SUBSTRING_INDEX(series.name, '@', 1) = SUBSTRING_INDEX(original_series.name, '@', 1) /* prefixes are equal */
		ORDER BY COALESCE(geographies.list_order, 999), geographies.handle`, seriesId)
	if err != nil {
		return
	}
	for rows.Next() {
		geography := models.Geography{}
		err = rows.Scan(
			&geography.FIPS,
			&geography.Handle,
			&geography.Name,
			&geography.ShortName,
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
