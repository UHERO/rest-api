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
		`SELECT
		  geographies.fips, geographies.display_name_short,
		  catgeo.chandle AS handle
		FROM
		  (SELECT DISTINCT(SUBSTRING_INDEX(SUBSTR(catnames.name, LOCATE('@', catnames.name) + 1), '.', 1)) AS chandle
		    FROM
		      (SELECT name
		        FROM series
		        WHERE (SELECT list FROM data_lists JOIN categories WHERE categories.data_list_id = data_lists.id AND (categories.id = ? OR categories.ancestry REGEXP CONCAT('[[:<:]]', ?, '[[:>:]]')))
		        REGEXP CONCAT('[[:<:]]', left(name, locate("@", name)), '.*[[:>:]]')) AS catnames) AS catgeo
		        LEFT JOIN geographies ON catgeo.chandle LIKE geographies.handle;`,
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
