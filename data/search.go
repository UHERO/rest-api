package data

import (
	"github.com/UHERO/rest-api/models"
)

func (r *SeriesRepository) GetSeriesBySearchText(searchText string) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT series.id, name, description, frequency, seasonally_adjusted,
	measurements.units_label, measurements.units_label_short, measurements.data_portal_name, measurements.percent, measurements.real,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM series LEFT JOIN geographies ON name LIKE CONCAT('%@', handle, '.%')
	LEFT JOIN measurements ON measurements.id = series.measurement_id
	WHERE
	((MATCH(name, description, dataPortalName)
	AGAINST(? IN NATURAL LANGUAGE MODE))
	OR LOWER(CONCAT(name, description, dataPortalName)) LIKE CONCAT('%', LOWER(?), '%')) LIMIT 50;`, searchText, searchText)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSearchSummary(searchText string) (searchSummary models.SearchSummary, err error) {
	searchSummary.SearchText = searchText
	rows, err := r.DB.Query(`SELECT MAX(SUBSTRING_INDEX(SUBSTR(name, LOCATE('@', name) + 1), '.', 1)) as geo,
  MAX(RIGHT(name, 1)) as freq
	FROM (SELECT name FROM series WHERE

    (MATCH(name, description, dataPortalName)
    AGAINST(? IN NATURAL LANGUAGE MODE))
     OR LOWER(CONCAT(name, description, dataPortalName)) LIKE CONCAT('%', LOWER(?), '%')
       ) AS s
    GROUP BY SUBSTR(name, LOCATE('@', name) + 1) ORDER BY COUNT(*) DESC;`, searchText, searchText)
	if err != nil {
		return
	}
	var geoFreq models.GeoFreq

	geoFreqs := map[string][]string{}
	freqGeos := map[string][]string{}

	for rows.Next() {
		err = rows.Scan(
			&geoFreq.Geography,
			&geoFreq.Frequency,
		)
		// set the default
		if searchSummary.DefaultGeoFreq == nil {
			searchSummary.DefaultGeoFreq = &models.GeoFreq{
				Geography: geoFreq.Geography,
				Frequency: geoFreq.Frequency,
			}
		}
		// add to the geoFreqs and freqGeos maps
		geoFreqs[geoFreq.Geography] = append(geoFreqs[geoFreq.Geography], geoFreq.Frequency)
		freqGeos[geoFreq.Frequency] = append(freqGeos[geoFreq.Frequency], geoFreq.Geography)
	}
	searchSummary.GeoFreqs = geoFreqs
	searchSummary.FreqGeos = freqGeos
	return
}

func (r *SeriesRepository) GetSearchResultsByGeoAndFreq(searchText string, geo string, freq string) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT series.id, name, description, frequency, seasonally_adjusted,
	measurements.units_label, measurements.units_label_short, measurements.data_portal_name, measurements.percent, measurements.real,
	fips, ?, display_name_short
	FROM series
	LEFT JOIN geographies ON handle LIKE ?
	LEFT JOIN measurements ON measurements.id = series.measurement_id
	WHERE
	((MATCH(name, description, dataPortalName)
	AGAINST(? IN NATURAL LANGUAGE MODE))
	OR LOWER(CONCAT(name, description, dataPortalName)) LIKE CONCAT('%', LOWER(?), '%'))
	AND LOWER(name) LIKE CONCAT('%@', LOWER(?), '.', LOWER(?)) LIMIT 50;`, geo, geo, searchText, searchText, geo, freq)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetInflatedSearchResultsByGeoAndFreq(
	searchText string,
	geo string,
	freq string,
) (seriesList []models.InflatedSeries, err error) {
	rows, err := r.DB.Query(`SELECT series.id, name, description, frequency, seasonally_adjusted,
	measurements.units_label, measurements.units_label_short, measurements.data_portal_name, measurements.percent, measurements.real,
	fips, ?, display_name_short
	FROM series
	LEFT JOIN geographies ON handle LIKE ?
	LEFT JOIN measurements ON measurements.id = series.measurement_id
	WHERE
	((MATCH(name, description, dataPortalName)
	AGAINST(? IN NATURAL LANGUAGE MODE))
	OR LOWER(CONCAT(name, description, dataPortalName)) LIKE CONCAT('%', LOWER(?), '%'))
	AND LOWER(name) LIKE CONCAT('%@', LOWER(?), '.', LOWER(?)) LIMIT 50;`, geo, geo, searchText, searchText, geo, freq)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		seriesObservations, scanErr := r.GetSeriesObservations(dataPortalSeries.Id)
		if scanErr != nil {
			return seriesList, scanErr
		}
		inflatedSeries := models.InflatedSeries{dataPortalSeries, seriesObservations}
		seriesList = append(seriesList, inflatedSeries)
	}
	return
}
