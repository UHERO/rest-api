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
	rows, err := r.DB.Query(`SELECT id, name, ancestry, default_handle, default_freq
							 FROM categories
							 WHERE universe = 'UHERO'
							 AND NOT hidden
							 ORDER BY categories.list_order;`)
	if err != nil {
		return
	}
	for rows.Next() {
		category := models.CategoryWithAncestry{}
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.Ancestry,
			&category.DefaultHandle,
			&category.DefaultFrequency,
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
	rows, err := r.DB.Query("SELECT id, name FROM categories WHERE ancestry IS NULL AND NOT hidden ORDER BY `list_order`;")
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
	var category models.CategoryWithAncestry
	err := r.DB.QueryRow(`SELECT
	categories.id, MAX(categories.name), MAX(ancestry),
	MIN(public_data_points.date) AS start_date, MAX(public_data_points.date) AS end_date
	FROM categories
	LEFT JOIN data_list_measurements ON categories.data_list_id = data_list_measurements.data_list_id
	LEFT JOIN measurement_series ON measurement_series.measurement_id = data_list_measurements.measurement_id
 	LEFT JOIN series ON series.id = measurement_series.series_id
 	LEFT JOIN public_data_points ON public_data_points.series_id = series.id
	WHERE categories.id = ?
	AND NOT categories.hidden
	AND NOT series.restricted
	GROUP BY categories.id;`, id).Scan(
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

	rows, err := r.DB.Query(
		`SELECT DISTINCT geo.fips, geo.display_name_short, geo.handle AS geo, RIGHT(series.name, 1) as freq
		FROM categories
		LEFT JOIN data_list_measurements ON data_list_measurements.data_list_id = categories.data_list_id
		LEFT JOIN measurement_series ON measurement_series.measurement_id = data_list_measurements.measurement_id
		LEFT JOIN series ON series.id = measurement_series.series_id
		JOIN geographies geo ON geo.id = series.geography_id
		WHERE categories.id = ?
		AND NOT categories.hidden
		AND NOT series.restricted;`, id)
	if err != nil {
		return dataPortalCategory, err
	}
	geoFreqs := map[string][]models.FrequencyResult{}
	geoByHandle := map[string]models.DataPortalGeography{}
	freqGeos := map[string][]models.DataPortalGeography{}
	freqByHandle := map[string]models.FrequencyResult{}
	for rows.Next() {
		scangeo := models.Geography{}
		frequency := models.FrequencyResult{}
		err = rows.Scan(
			&scangeo.FIPS,
			&scangeo.Name,
			&scangeo.Handle,
			&frequency.Freq,
		)
		geography := models.DataPortalGeography{Handle: scangeo.Handle}
		if scangeo.FIPS.Valid {
			geography.FIPS = scangeo.FIPS.String
		}
		if scangeo.Name.Valid {
			geography.Name = scangeo.Name.String
		}
		frequency.Label = freqLabel[frequency.Freq]
		// update the freq and geo maps
		geoByHandle[geography.Handle] = geography
		freqByHandle[frequency.Freq] = frequency
		// add to the geoFreqs and freqGeos maps
		geoFreqs[geography.Handle] = append(geoFreqs[geography.Handle], frequency)
		freqGeos[frequency.Freq] = append(freqGeos[frequency.Freq], geography)
	}
	geoFreqsResult := []models.GeographyFrequencies{}
	for geo, freqs := range geoFreqs {
		sort.Sort(models.ByFrequency(freqs))
		geoFreqsResult = append(geoFreqsResult, models.GeographyFrequencies{
			DataPortalGeography: geoByHandle[geo],
			Frequencies:         freqs,
		})
	}

	freqGeosResult := []models.FrequencyGeographies{}
	for _, freq := range models.FreqOrder {
		if val, ok := freqByHandle[freq]; ok {
			freqGeosResult = append(freqGeosResult, models.FrequencyGeographies{
				FrequencyResult: val,
				Geographies:     freqGeos[freq],
			})
		}
	}

	dataPortalCategory.GeographyFrequencies = &geoFreqsResult
	dataPortalCategory.FrequencyGeographies = &freqGeosResult
	dataPortalCategory.GeoFreqsDeprecated = &geoFreqsResult
	dataPortalCategory.FreqGeosDeprecated = &freqGeosResult
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
		category := models.CategoryWithAncestry{}
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
