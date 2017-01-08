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
	WHERE NOT restricted AND
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
	rows, err := r.DB.Query(`SELECT geographies.fips, geographies.display_name_short, geofreq.geo, geofreq.freq
FROM (SELECT MAX(SUBSTRING_INDEX(SUBSTR(name, LOCATE('@', name) + 1), '.', 1)) as geo,
       MAX(RIGHT(name, 1)) as freq
FROM (SELECT name FROM series WHERE NOT restricted AND
                                    (MATCH(name, description, dataPortalName)
                                    AGAINST(? IN NATURAL LANGUAGE MODE))
                                    OR LOWER(CONCAT(name, description, dataPortalName)) LIKE CONCAT('%', LOWER(?), '%')) AS s
GROUP BY SUBSTR(name, LOCATE('@', name) + 1) ORDER BY COUNT(*) DESC) as geofreq
LEFT JOIN geographies ON geographies.handle = geofreq.geo;`, searchText, searchText)
	if err != nil {
		return
	}
	legacyGeoFreqs := map[string][]string{}
	legacyFreqGeos := map[string][]string{}
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
		// set the default as the first in the sorted list
		if searchSummary.DefaultGeoFreq == nil {
			searchSummary.DefaultGeoFreq = &models.GeographyFrequency{
				Geography: geography,
				Frequency: frequency,
			}
		}
		// update the freq and geo maps
		geoByHandle[geography.Handle] = geography
		freqByHandle[frequency.Freq] = frequency
		// add to the geoFreqs and freqGeos maps
		geoFreqs[geography.Handle] = append(geoFreqs[geography.Handle], frequency)
		freqGeos[frequency.Freq] = append(freqGeos[frequency.Freq], geography)

		legacyGeoFreqs[geography.Handle] = append(legacyGeoFreqs[geography.Handle], frequency.Freq)
		legacyFreqGeos[frequency.Freq] = append(legacyFreqGeos[frequency.Freq], geography.Handle)
	}

	geoFreqsResult := []models.GeographyFrequencies{}
	for geo, freqs := range geoFreqs {
		geoFreqsResult = append(geoFreqsResult, models.GeographyFrequencies{
			DataPortalGeography: geoByHandle[geo],
			Frequencies: freqs,
		})
	}

	freqGeosResult := []models.FrequencyGeographies{}
	for freq, geos := range freqGeos {
		freqGeosResult = append(freqGeosResult, models.FrequencyGeographies{
			FrequencyResult: freqByHandle[freq],
			Geographies: geos,
		})
	}

	searchSummary.GeoFreqs = legacyGeoFreqs
	searchSummary.FreqGeos = legacyFreqGeos
	searchSummary.GeographyFrequencies = &geoFreqsResult
	searchSummary.FrequencyGeographies = &freqGeosResult
	return
}

func (r *SeriesRepository) GetSearchResultsByGeoAndFreq(searchText string, geo string, freq string) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT series.id, name, description, frequency, seasonally_adjusted,
	measurements.units_label, measurements.units_label_short, measurements.data_portal_name, measurements.percent, measurements.real,
	fips, ?, display_name_short
	FROM series
	LEFT JOIN geographies ON handle LIKE ?
	LEFT JOIN measurements ON measurements.id = series.measurement_id
	WHERE NOT restricted AND
	((MATCH(name, description, dataPortalName)
	AGAINST(? IN NATURAL LANGUAGE MODE))
	OR LOWER(CONCAT(name, description, dataPortalName)) LIKE CONCAT('%', LOWER(?), '%'))
	AND LOWER(name) LIKE CONCAT('%@', LOWER(?), '.', LOWER(?)) LIMIT 50;`,
		geo,
		geo,
		searchText,
		searchText,
		geo,
		freq,
	)
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
	WHERE NOT series.restricted AND
	((MATCH(name, description, dataPortalName)
	AGAINST(? IN NATURAL LANGUAGE MODE))
	OR LOWER(CONCAT(name, description, dataPortalName)) LIKE CONCAT('%', LOWER(?), '%'))
	AND LOWER(name) LIKE CONCAT('%@', LOWER(?), '.', LOWER(?)) LIMIT 50;`,
		geo,
		geo,
		searchText,
		searchText,
		geo,
		freq,
	)
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
