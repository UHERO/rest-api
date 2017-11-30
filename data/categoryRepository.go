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
			ANY_VALUE(categories.name) AS catname,
			ANY_VALUE(categories.ancestry) AS ancest,
			categories.default_freq AS catfreq,
			ANY_VALUE(geographies.handle) AS catgeo,
			ANY_VALUE(geographies.fips) AS catgeofips,
			ANY_VALUE(geographies.display_name) AS catgeoname,
			ANY_VALUE(geographies.display_name_short) AS catgeonameshort,
			MIN(public_data_points.date) AS startdate,
			MAX(public_data_points.date) AS enddate
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
		category := models.CategoryWithAncestry{}
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.Ancestry,
			&category.DefaultFrequency,
			&category.DefaultGeoHandle,
			&category.DefaultGeoFIPS,
			&category.DefaultGeoName,
			&category.DefaultGeoShortName,
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
		if category.DefaultFrequency.Valid || category.DefaultGeoHandle.Valid || category.ObservationStart.Valid || category.ObservationEnd.Valid {
			// Only initialize Defaults struct if any defaults values are available
			dataPortalCategory.Defaults = &models.CategoryDefaults{}
		}
		if category.DefaultFrequency.Valid {
			dataPortalCategory.Defaults.Frequency = &models.DataPortalFrequency{
				Freq: category.DefaultFrequency.String,
				Label: freqLabel[category.DefaultFrequency.String],
			}
		}
		if category.DefaultGeoHandle.Valid {
			dataPortalCategory.Defaults.Geography = &models.DataPortalGeography{
				Handle: category.DefaultGeoHandle.String,
			}
			if category.DefaultGeoFIPS.Valid {
				dataPortalCategory.Defaults.Geography.FIPS = category.DefaultGeoFIPS.String
			}
			if category.DefaultGeoName.Valid {
				dataPortalCategory.Defaults.Geography.Name = category.DefaultGeoName.String
			}
			if category.DefaultGeoShortName.Valid {
				dataPortalCategory.Defaults.Geography.ShortName = category.DefaultGeoShortName.String
			}
		}
		if category.ObservationStart.Valid {
			dataPortalCategory.Defaults.ObservationStart = &category.ObservationStart.Time
		}
		if category.ObservationEnd.Valid {
			dataPortalCategory.Defaults.ObservationEnd = &category.ObservationEnd.Time
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
	return r.GetCategoryByIdGeoFreq(id, "", "")
}

func (r *CategoryRepository) GetCategoryByIdGeoFreq(id int64, originGeo string, originFreq string) (models.Category, error) {
	dataPortalCategory := models.Category{}
	rows, err := r.DB.Query(
		`SELECT categories.id,
			ANY_VALUE(categories.name) AS catname,
			ANY_VALUE(parentcat.id) AS parent_id,
			ANY_VALUE(COALESCE(mygeo.handle, parentgeo.handle)) AS def_geo,
			ANY_VALUE(COALESCE(categories.default_freq, parentcat.default_freq)) AS def_freq,
			ANY_VALUE(geographies.handle) AS series_geo,
			RIGHT(series.name, 1) AS series_freq,
			ANY_VALUE(geographies.fips) AS geofips,
			ANY_VALUE(geographies.display_name) AS geoname,
			ANY_VALUE(geographies.display_name_short) AS geonameshort,
			MIN(public_data_points.date) AS startdate,
			MAX(public_data_points.date) AS enddate
		FROM categories
  		LEFT JOIN geographies mygeo ON mygeo.id = categories.default_geo_id
		LEFT JOIN categories parentcat ON parentcat.id = SUBSTRING_INDEX(categories.ancestry, '/', -1)
		LEFT JOIN geographies parentgeo ON parentgeo.id = parentcat.default_geo_id
	        LEFT JOIN data_list_measurements ON data_list_measurements.data_list_id = categories.data_list_id
		LEFT JOIN measurement_series ON measurement_series.measurement_id = data_list_measurements.measurement_id
		LEFT JOIN series
		    ON series.id = measurement_series.series_id
		   AND NOT series.restricted
		LEFT JOIN geographies ON geographies.id = series.geography_id
		LEFT JOIN public_data_points ON public_data_points.series_id = series.id
		WHERE categories.id = ?
		AND public_data_points.value IS NOT null /* suppress those with no public data */
		GROUP BY categories.id, geographies.id, RIGHT(series.name, 1)  ;`, id)
	if err != nil {
		return dataPortalCategory, err
	}
	var geosResult  []models.DataPortalGeography
	var freqsResult []models.DataPortalFrequency
	//var defaultGeo	*models.DataPortalGeography
	//var defaultFreq *models.DataPortalFrequency
	//seenGeos := map[string]*models.DataPortalGeography{}
	//seenFreqs := map[string]*models.DataPortalFrequency{}

	for rows.Next() {
		//var isDefaultGeo, isDefaultFreq	sql.NullBool
		var handle, seriesFreq string
		category := models.CategoryWithAncestry{}
		scangeo := models.Geography{}
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.ParentId,
			&category.DefaultGeoHandle,
			&category.DefaultFrequency,
			&handle,
			&seriesFreq,
			&scangeo.FIPS,
			&scangeo.Name,
			&scangeo.ShortName,
			&scangeo.ObservationStart,
			&scangeo.ObservationEnd,
		)
		if dataPortalCategory.Id == 0 {
			// Store Category-level information
			dataPortalCategory.Id = category.Id
			dataPortalCategory.Name = category.Name
			if category.ParentId.Valid {
				dataPortalCategory.ParentId = category.ParentId.Int64
			}
			if originGeo == "" && category.DefaultGeoHandle.Valid {
				originGeo = category.DefaultGeoHandle.String
			}
			if originFreq == "" && category.DefaultFrequency.Valid {
				originFreq = category.DefaultFrequency.String
			}
			startOfTime := time.Date(1099, 1, 1, 0, 0, 0, 0, &time.Location{})
			endOfTime := time.Date(2999, 1, 1, 0, 0, 0, 0, &time.Location{})
			dataPortalCategory.ObservationStart = &endOfTime // yes, reversed at init time
			dataPortalCategory.ObservationEnd = &startOfTime
		}
		geo := &models.DataPortalGeography{Handle: handle}
		freq := &models.DataPortalFrequency{Freq: seriesFreq, Label: freqLabel[seriesFreq]}
		if scangeo.FIPS.Valid {
			geo.FIPS = scangeo.FIPS.String
		}
		if scangeo.Name.Valid {
			geo.Name = scangeo.Name.String
		}
		if scangeo.ShortName.Valid {
			geo.ShortName = scangeo.ShortName.String
		}
		if scangeo.ObservationStart.Valid  {
			 geo.ObservationStart = &scangeo.ObservationStart.Time
			freq.ObservationStart = &scangeo.ObservationStart.Time
			if scangeo.ObservationStart.Time.Before(*dataPortalCategory.ObservationStart) {
				dataPortalCategory.ObservationStart = &scangeo.ObservationStart.Time
			}
		}
		if scangeo.ObservationEnd.Valid  {
			 geo.ObservationEnd = &scangeo.ObservationEnd.Time
			freq.ObservationEnd = &scangeo.ObservationEnd.Time
			if scangeo.ObservationEnd.Time.After(*dataPortalCategory.ObservationEnd) {
				dataPortalCategory.ObservationEnd = &scangeo.ObservationEnd.Time
			}
		}
		if seriesFreq == originFreq {
			geosResult = append(geosResult, *geo)

			if geo.Handle == originGeo {
				dataPortalCategory.Current = &models.CurrentGeoFreq{
					Geography: originGeo,
					Frequency: originFreq,
					ObservationStart: &scangeo.ObservationStart.Time,
					ObservationEnd: &scangeo.ObservationEnd.Time,
				}
			}
		}
		if geo.Handle == originGeo {
			freqsResult = append(freqsResult, *freq)
		}
	}
	sort.Sort(models.ByGeography(geosResult))
	sort.Sort(models.ByFrequency(freqsResult))
	dataPortalCategory.Geographies = &geosResult
	dataPortalCategory.Frequencies = &freqsResult
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
