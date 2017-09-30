package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CategoryRepository struct {
	DB *sql.DB
}

func (r *CategoryRepository) GetAllCategories() (categories []models.Category, err error) {
	categories, err = r.GetAllCategoriesByUniverse("UHERO")
	return
}

func (r *CategoryRepository) GetAllCategoriesByUniverse(universe string) (categories []models.Category, err error) {
	rows, err := r.DB.Query(
		`SELECT categories.id,
			categories.default_freq AS catfreq,
			ANY_VALUE(categories.name) AS catname,
			ANY_VALUE(categories.ancestry) AS ancest,
			ANY_VALUE(geographies.handle) AS catgeo,
			ANY_VALUE(geographies.fips) AS catgeofips,
			ANY_VALUE(geographies.display_name_short) AS catgeodisp,
			MIN(public_data_points.date) as startdate, MAX(public_data_points.date) as enddate
		FROM categories
		LEFT JOIN geographies ON geographies.id = categories.default_geo_id
		LEFT JOIN data_list_measurements ON data_list_measurements.data_list_id = categories.data_list_id
		LEFT JOIN measurement_series ON measurement_series.measurement_id = data_list_measurements.measurement_id
		LEFT JOIN series
		    ON series.id = measurement_series.series_id
		   AND series.geography_id = categories.default_geo_id
		   AND RIGHT(series.name, 1) = categories.default_freq
		   AND NOT series.restricted
		LEFT JOIN public_data_points ON public_data_points.series_id = series.id
		WHERE categories.universe = ?
		AND categories.ancestry IS NOT NULL
		AND NOT categories.hidden
		GROUP BY categories.id, categories.default_geo_id, categories.default_freq
		ORDER BY categories.list_order;`, universe)
	if err != nil {
		return
	}
	for rows.Next() {
		category := models.CategoryWithAncestryEtc{}
		err = rows.Scan(
			&category.Id,
			&category.DefaultFrequency,
			&category.Name,
			&category.Ancestry,
			&category.DefaultGeoHandle,
			&category.DefaultGeoFIPS,
			&category.DefaultGeoName,
			&category.ObservationStart,
			&category.ObservationEnd,
		)
		if err != nil {
			return
		}

		parentId := getParentId(category.Ancestry)
		dataPortalCategory := models.Category{
			Id:       category.Id,
			Name:     category.Name,
			ParentId: parentId,
		}
		if category.DefaultHandle.Valid && category.DefaultFrequency.Valid {
			dataPortalCategory.DefaultGeoFreq = &models.GeoFreq{
				Geography: category.DefaultHandle.String,
				Frequency: category.DefaultFrequency.String,
			}
		}
		categories = append(categories, dataPortalCategory)
	}
	return
}

func getParentId(ancestry sql.NullString) (parentId int64) {
	if !ancestry.Valid {
		return
	}
	parents := strings.Split(ancestry.String, "/")
	if len(parents) == 0 {
		return
	}
	parentId, _ = strconv.ParseInt(parents[len(parents)-1], 10, 64)
	return
}

func (r *CategoryRepository) GetCategoryRoots() (categories []models.Category, err error) {
	rows, err := r.DB.Query(`SELECT id, name FROM categories
				WHERE ancestry IS NULL
				AND NOT hidden
				ORDER BY list_order;`)
	if err != nil {
		return
	}
	for rows.Next() {
		category := models.Category{}
		err = rows.Scan(
			&category.Id,
			&category.Name,
		)
		if err != nil {
			return
		}
		categories = append(categories, category)
	}
	return
}

func (r *CategoryRepository) GetCategoryById(id int64) (models.Category, error) {
	var category models.CategoryWithAncestryEtc
	err := r.DB.QueryRow(
		`SELECT categories.id, ANY_VALUE(categories.name), ANY_VALUE(ancestry),
		MIN(public_data_points.date), MAX(public_data_points.date)
		FROM categories
		LEFT JOIN data_list_measurements ON categories.data_list_id = data_list_measurements.data_list_id
		LEFT JOIN measurement_series ON measurement_series.measurement_id = data_list_measurements.measurement_id
		LEFT JOIN series ON series.id = measurement_series.series_id
		LEFT JOIN public_data_points ON public_data_points.series_id = series.id
		WHERE categories.id = ?
		AND NOT categories.hidden
		AND NOT series.restricted
		GROUP BY categories.id ;`, id).Scan(
		&category.Id,
		&category.Name,
		&category.Ancestry,
		&category.ObservationStart,
		&category.ObservationEnd,
	)
	parentId := getParentId(category.Ancestry)
	dataPortalCategory := models.Category{
		Id:       category.Id,
		Name:     category.Name,
		ParentId: parentId,
	}
	if category.ObservationStart.Valid && category.ObservationStart.Time.After(time.Time{}) {
		dataPortalCategory.ObservationStart = &category.ObservationStart.Time
	}
	if category.ObservationEnd.Valid && category.ObservationEnd.Time.After(time.Time{}) {
		dataPortalCategory.ObservationEnd = &category.ObservationEnd.Time
	}

	rows, err := r.DB.Query(`SELECT ANY_VALUE(geographies.fips), ANY_VALUE(geographies.display_name_short),
					ANY_VALUE(geographies.handle), ANY_VALUE(RIGHT(series.name, 1)),
					MIN(public_data_points.date), MAX(public_data_points.date)
			FROM categories
			LEFT JOIN data_list_measurements ON data_list_measurements.data_list_id = categories.data_list_id
			LEFT JOIN measurement_series ON measurement_series.measurement_id = data_list_measurements.measurement_id
			LEFT JOIN series ON series.id = measurement_series.series_id
			LEFT JOIN public_data_points ON public_data_points.series_id = series.id
			LEFT JOIN geographies ON geographies.id = series.geography_id
			WHERE categories.id = ?
			AND NOT categories.hidden
			AND NOT series.restricted
			GROUP BY geographies.id, RIGHT(series.name, 1) ORDER BY COUNT(*) DESC;`, id)
	if err != nil {
		return dataPortalCategory, err
	}
	geoFreqs := map[string][]models.DataPortalFrequency{}
	geoByHandle := map[string]models.DataPortalGeography{}

	for rows.Next() {
		scangeo := models.Geography{}
		frequency := models.DataPortalFrequency{}
		err = rows.Scan(
			&scangeo.FIPS,
			&scangeo.Name,
			&scangeo.Handle,
			&frequency.Freq,
			&scangeo.ObservationStart,
			&scangeo.ObservationEnd,
		)
		geography := models.DataPortalGeography{Handle: scangeo.Handle}
		if scangeo.ObservationStart.Valid && scangeo.ObservationStart.Time.After(time.Time{}) {
			geography.ObservationStart = &scangeo.ObservationStart.Time
			frequency.ObservationStart = geography.ObservationStart
		}
		if scangeo.ObservationEnd.Valid && scangeo.ObservationEnd.Time.After(time.Time{}) {
			geography.ObservationEnd = &scangeo.ObservationEnd.Time
			frequency.ObservationEnd = geography.ObservationEnd
		}
		if scangeo.FIPS.Valid {
			geography.FIPS = scangeo.FIPS.String
		}
		if scangeo.Name.Valid {
			geography.Name = scangeo.Name.String
		}
		frequency.Label = freqLabel[frequency.Freq]
		geoByHandle[geography.Handle] = geography
		geoFreqs[geography.Handle] = append(geoFreqs[geography.Handle], frequency)
	}
	geoFreqsResult := []models.GeographyFrequencies{}
	for geo, freqs := range geoFreqs {
		sort.Sort(models.ByFrequency(freqs))
		geoFreqsResult = append(geoFreqsResult, models.GeographyFrequencies{
			DataPortalGeography: geoByHandle[geo],
			Frequencies:         freqs,
		})
	}
	dataPortalCategory.GeographyFrequencies = &geoFreqsResult
	return dataPortalCategory, err
}

func (r *CategoryRepository) GetCategoriesByName(name string) (categories []models.Category, err error) {
	fuzzyString := "%" + name + "%"
	rows, err := r.DB.Query(`SELECT id, name, ancestry FROM categories
							 WHERE LOWER(name) LIKE ? AND NOT hidden
							 ORDER BY list_order;`, fuzzyString)
	if err != nil {
		return
	}
	for rows.Next() {
		category := models.CategoryWithAncestryEtc{}
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.Ancestry,
		)
		if err != nil {
			return
		}
		parentId := getParentId(category.Ancestry)
		categories = append(categories, models.Category{
			Id:       category.Id,
			Name:     category.Name,
			ParentId: parentId,
		})
	}
	return
}

func (r *CategoryRepository) GetChildrenOf(id int64) (categories []models.Category, err error) {
	rows, err := r.DB.Query(`SELECT id, name, parent_id FROM categories
							 WHERE parent_id = ? AND NOT hidden
							 ORDER BY list_order;`, id)
	if err != nil {
		return
	}
	for rows.Next() {
		category := models.Category{}
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.ParentId,
		)
		if err != nil {
			return
		}
		categories = append(categories, category)
	}
	return
}
