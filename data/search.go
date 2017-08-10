package data

import (
	"sort"
	"time"

	"github.com/UHERO/rest-api/models"
)

var freqNames map[string]string = map[string]string{
	"A": "year",
	"Q": "quarter",
	"M": "month",
	"S": "semi",
	"W": "week",
	"D": "day",
}

func (r *SeriesRepository) GetSeriesBySearchTextAndUniverse(searchText string, universeText string) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT
	series.id, series.name, series.description, frequency, series.seasonally_adjusted, series.seasonal_adjustment,
	COALESCE(NULLIF(units.long_label, ''), NULLIF(MAX(measurement_units.long_label), '')),
	COALESCE(NULLIF(units.short_label, ''), NULLIF(MAX(measurement_units.short_label), '')),
	COALESCE(NULLIF(series.dataPortalName, ''), MAX(measurements.data_portal_name)), series.percent, series.real,
	COALESCE(NULLIF(sources.description, ''), NULLIF(MAX(measurement_sources.description), '')),
	COALESCE(NULLIF(series.source_link, ''), NULLIF(MAX(measurements.source_link), ''), NULLIF(sources.link, ''), NULLIF(MAX(measurement_sources.link), '')),
	COALESCE(NULLIF(source_details.description, ''), NULLIF(MAX(measurement_source_details.description), '')),
	MAX(measurements.table_prefix), MAX(measurements.table_postfix),
	MAX(measurements.id), MAX(measurements.data_portal_name),
	NULL, series.base_year, series.decimals,
	MAX(geo.fips), MAX(geo.handle) AS shandle, MAX(geo.display_name_short)
	FROM series
	JOIN geographies geo ON geo.id = series.geography_id
	LEFT JOIN measurement_series ON measurement_series.series_id = series.id
	LEFT JOIN measurements ON measurements.id = measurement_series.measurement_id
	LEFT JOIN units ON units.id = series.unit_id
	LEFT JOIN units AS measurement_units ON measurement_units.id = measurements.unit_id
	LEFT JOIN sources ON sources.id = series.source_id
	LEFT JOIN sources AS measurement_sources ON measurement_sources.id = measurements.source_id
	LEFT JOIN source_details ON source_details.id = series.source_detail_id
	LEFT JOIN source_details AS measurement_source_details ON measurement_source_details.id = measurements.source_detail_id
	LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE series.universe = UPPER(?)
	AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)
	AND ((MATCH(series.name, series.description, series.dataPortalName) AGAINST(? IN NATURAL LANGUAGE MODE))
	  OR LOWER(CONCAT(series.name, series.description, series.dataPortalName)) LIKE CONCAT('%', LOWER(?), '%'))
	GROUP BY series.id
	LIMIT 50;`, universeText, searchText, searchText)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geoFreqs, freqGeos, err := getFreqGeoCombinations(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.GeographyFrequencies = &geoFreqs
		dataPortalSeries.FrequencyGeographies = &freqGeos
		dataPortalSeries.GeoFreqsDeprecated = &geoFreqs
		dataPortalSeries.FreqGeosDeprecated = &freqGeos
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesBySearchText(searchText string) (seriesList []models.DataPortalSeries, err error) {
	seriesList, err = r.GetSeriesBySearchTextAndUniverse(searchText, "UHERO")
	return
}

func (r *SeriesRepository) GetSearchSummaryByUniverse(searchText string, universeText string) (searchSummary models.SearchSummary, err error) {
	searchSummary.SearchText = searchText

	var observationStart, observationEnd models.NullTime
	err = r.DB.QueryRow(`SELECT MIN(public_data_points.date) AS start_date, MAX(public_data_points.date) AS end_date
	FROM series
 	LEFT JOIN public_data_points ON public_data_points.series_id = series.id
 	LEFT JOIN feature_toggles ON feature_toggles.universe = public_data_points.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE public_data_points.universe = UPPER(?)
	AND ((MATCH(series.name, series.description, series.dataPortalName) AGAINST(? IN NATURAL LANGUAGE MODE))
	  OR LOWER(CONCAT(series.name, series.description, series.dataPortalName)) LIKE CONCAT('%', LOWER(?), '%'))
	AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)`, universeText, searchText, searchText).Scan(
		&observationStart,
		&observationEnd,
	)
	if err != nil {
		return
	}
	if observationStart.Valid && observationStart.Time.After(time.Time{}) {
		searchSummary.ObservationStart = &observationStart.Time
	}
	if observationEnd.Valid && observationEnd.Time.After(time.Time{}) {
		searchSummary.ObservationEnd = &observationEnd.Time
	}

	rows, err := r.DB.Query(`
	SELECT DISTINCT g.fips, g.display_name_short, g.handle AS geo, RIGHT(series.name, 1) as freq
	FROM series
	  JOIN geographies g on g.id = series.geography_id
    	  LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE series.universe = UPPER(?)
	AND NOT restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT quarantined)
	AND ((MATCH(series.name, series.description, dataPortalName) AGAINST(? IN NATURAL LANGUAGE MODE))
	   OR LOWER(CONCAT(series.name, series.description, dataPortalName)) LIKE CONCAT('%', LOWER(?), '%'))
	ORDER BY 1,2,3,4;`, universeText, searchText, searchText)
	if err != nil {
		return
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

	searchSummary.GeographyFrequencies = &geoFreqsResult
	searchSummary.FrequencyGeographies = &freqGeosResult
	searchSummary.GeoFreqsDeprecated = &geoFreqsResult
	searchSummary.FreqGeosDeprecated = &freqGeosResult
	return
}

func (r *SeriesRepository) GetSearchSummary(searchText string) (searchSummary models.SearchSummary, err error) {
	searchSummary, err = r.GetSearchSummaryByUniverse(searchText, "UHERO")
	return
}

func (r *SeriesRepository) GetSearchResultsByGeoAndFreqAndUniverse(
	searchText string,
	geo string,
	freq string,
	universeText string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT
	series.id, series.name, series.description, frequency, series.seasonally_adjusted, series.seasonal_adjustment,
	COALESCE(NULLIF(units.long_label, ''), NULLIF(MAX(measurement_units.long_label), '')),
	COALESCE(NULLIF(units.short_label, ''), NULLIF(MAX(measurement_units.short_label), '')),
	COALESCE(NULLIF(series.dataPortalName, ''), MAX(measurements.data_portal_name)), series.percent, series.real,
	COALESCE(NULLIF(sources.description, ''), NULLIF(MAX(measurement_sources.description), '')),
	COALESCE(NULLIF(series.source_link, ''), NULLIF(MAX(measurements.source_link), ''), NULLIF(sources.link, ''), NULLIF(MAX(measurement_sources.link), '')),
	COALESCE(NULLIF(source_details.description, ''), NULLIF(MAX(measurement_source_details.description), '')),
	MAX(measurements.table_prefix), MAX(measurements.table_postfix),
	MAX(measurements.id), MAX(measurements.data_portal_name),
	NULL, series.base_year, series.decimals,
	MAX(geo.fips), MAX(geo.handle), MAX(geo.display_name_short)
	FROM series
	JOIN geographies geo ON geo.id = series.geography_id
	LEFT JOIN measurement_series ON measurement_series.series_id = series.id
	LEFT JOIN measurements ON measurements.id = measurement_series.measurement_id
	LEFT JOIN units ON units.id = series.unit_id
	LEFT JOIN units AS measurement_units ON measurement_units.id = measurements.unit_id
	LEFT JOIN sources ON sources.id = series.source_id
	LEFT JOIN sources AS measurement_sources ON measurement_sources.id = measurements.source_id
	LEFT JOIN source_details ON source_details.id = series.source_detail_id
	LEFT JOIN source_details AS measurement_source_details ON measurement_source_details.id = measurements.source_detail_id
	LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE series.universe = UPPER(?)
	AND geo.handle = ?
	AND series.frequency = ?
	AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)
	AND ((MATCH(series.name, series.description, series.dataPortalName) AGAINST(? IN NATURAL LANGUAGE MODE))
	  OR LOWER(CONCAT(series.name, series.description, series.dataPortalName)) LIKE CONCAT('%', LOWER(?), '%'))
	GROUP BY series.id
	LIMIT 50;`,
		universeText,
		geo,
		freqNames[freq],
		searchText,
		searchText,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geoFreqs, freqGeos, err := getFreqGeoCombinations(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.GeographyFrequencies = &geoFreqs
		dataPortalSeries.FrequencyGeographies = &freqGeos
		dataPortalSeries.GeoFreqsDeprecated = &geoFreqs
		dataPortalSeries.FreqGeosDeprecated = &freqGeos
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSearchResultsByGeoAndFreq(searchText string, geo string, freq string) (seriesList []models.DataPortalSeries, err error) {
	seriesList, err = r.GetSearchResultsByGeoAndFreqAndUniverse(searchText, geo, freq, "UHERO")
	return
}

func (r *SeriesRepository) GetInflatedSearchResultsByGeoAndFreqAndUniverse(
	searchText string,
	geo string,
	freq string,
	universeText string,
) (seriesList []models.InflatedSeries, err error) {
	rows, err := r.DB.Query(`SELECT DISTINCT
	series.id, series.name, series.description, frequency, series.seasonally_adjusted, series.seasonal_adjustment,
	COALESCE(NULLIF(units.long_label, ''), NULLIF(MAX(measurement_units.long_label), '')),
	COALESCE(NULLIF(units.short_label, ''), NULLIF(MAX(measurement_units.short_label), '')),
	COALESCE(NULLIF(series.dataPortalName, ''), MAX(measurements.data_portal_name)), series.percent, series.real,
	COALESCE(NULLIF(sources.description, ''), NULLIF(MAX(measurement_sources.description), '')),
	COALESCE(NULLIF(series.source_link, ''), NULLIF(MAX(measurements.source_link), ''), NULLIF(sources.link, ''), NULLIF(MAX(measurement_sources.link), '')),
	COALESCE(NULLIF(source_details.description, ''), NULLIF(MAX(measurement_source_details.description), '')),
	MAX(measurements.table_prefix), MAX(measurements.table_postfix),
	MAX(measurements.id), MAX(measurements.data_portal_name),
	NULL, series.base_year, series.decimals,
	MAX(geo.fips), MAX(geo.handle), MAX(geo.display_name_short)
	FROM series
	JOIN geographies geo ON geo.id = series.geography_id
	LEFT JOIN measurement_series ON measurement_series.series_id = series.id
	LEFT JOIN measurements ON measurements.id = measurement_series.measurement_id
	LEFT JOIN units ON units.id = series.unit_id
	LEFT JOIN units AS measurement_units ON measurement_units.id = measurements.unit_id
	LEFT JOIN sources ON sources.id = series.source_id
	LEFT JOIN sources AS measurement_sources ON measurement_sources.id = measurements.source_id
	LEFT JOIN source_details ON source_details.id = series.source_detail_id
	LEFT JOIN source_details AS measurement_source_details ON measurement_source_details.id = measurements.source_detail_id
	LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE series.universe = UPPER(?)
	AND geo.handle = ?
	AND series.frequency = ?
	AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)
	AND ((MATCH(series.name, series.description, series.dataPortalName) AGAINST(? IN NATURAL LANGUAGE MODE))
	  OR LOWER(CONCAT(series.name, series.description, series.dataPortalName)) LIKE CONCAT('%', LOWER(?), '%'))
	GROUP BY series.id
	LIMIT 50;`,
		universeText,
		geo,
		freqNames[freq],
		searchText,
		searchText,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geoFreqs, freqGeos, err := getFreqGeoCombinations(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.GeographyFrequencies = &geoFreqs
		dataPortalSeries.FrequencyGeographies = &freqGeos
		dataPortalSeries.GeoFreqsDeprecated = &geoFreqs
		dataPortalSeries.FreqGeosDeprecated = &freqGeos
		seriesObservations, scanErr := r.GetSeriesObservations(dataPortalSeries.Id)
		if scanErr != nil {
			return seriesList, scanErr
		}
		inflatedSeries := models.InflatedSeries{dataPortalSeries, seriesObservations}
		seriesList = append(seriesList, inflatedSeries)
	}
	return
}

func (r *SeriesRepository) GetInflatedSearchResultsByGeoAndFreq(searchText string, geo string, freq string) (seriesList []models.InflatedSeries, err error) {
	seriesList, err = r.GetInflatedSearchResultsByGeoAndFreqAndUniverse(searchText, geo, freq, "UHERO")
	return
}
