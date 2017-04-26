package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"strings"
	"time"
)

type SeriesRepository struct {
	DB *sql.DB
}

type transformation struct {
	Statement        string
	PlaceholderCount int
	Label            string
}

type GroupType int

const (
	Category GroupType = iota
	Measurement
)

const (
	Levels           = "lvl"
	YOYPercentChange = "pc1"
	YTDPercentChange = "ytdpc1"
	YOYChange        = "ch1"
	YTDChange        = "ytdch1"
)

var transformations map[string]transformation = map[string]transformation{
	Levels: { // untransformed value
		Statement: `SELECT date, value/units, (pseudo_history = b'1')
		FROM data_points
		LEFT JOIN series ON data_points.series_id = series.id
		WHERE series_id = ? and current = 1;`,
		PlaceholderCount: 1,
		Label:            "lvl",
	},
	YOYPercentChange: { // percent change from 1 year ago
		Statement: `SELECT t1.date, (t1.value/t2.last_value - 1)*100 AS yoy,
				(t1.pseudo_history = b'1') AND (t2.pseudo_history = b'1') AS ph
				FROM (SELECT value, date, pseudo_history, DATE_SUB(date, INTERVAL 1 YEAR) AS last_year
				      FROM data_points WHERE series_id = ? AND current = 1) AS t1
				LEFT JOIN (SELECT value AS last_value, date, pseudo_history
				           FROM data_points WHERE series_id = ? and current = 1) AS t2 ON (t1.last_year = t2.date);`,
		PlaceholderCount: 2,
		Label:            "pc1",
	},
	YOYChange: { // change from 1 year ago
		Statement: `SELECT t1.date, (t1.value - t2.last_value)/series.units AS yoy,
				(t1.pseudo_history = b'1') AND (t2.pseudo_history = b'1') AS ph
				FROM (SELECT series_id, value, date, pseudo_history, DATE_SUB(date, INTERVAL 1 YEAR) AS last_year
				      FROM data_points WHERE series_id = ? AND current = 1) AS t1
				LEFT JOIN (SELECT value AS last_value, date, pseudo_history
				           FROM data_points WHERE series_id = ? and current = 1) AS t2 ON (t1.last_year = t2.date)
				LEFT JOIN series ON t1.series_id = series.id;`,
		PlaceholderCount: 2,
		Label:            "pc1",
	},
	YTDChange: { // ytd change from 1 year ago
		Statement: `SELECT t1.date, (t1.ytd - t2.last_ytd)/series.units AS ytd,
				(t1.pseudo_history = b'1') AND (t2.pseudo_history = b'1') AS ph
      FROM (SELECT date, value, series_id, pseudo_history, @sum := IF(@year = YEAR(date), @sum, 0) + value AS ytd,
              @year := year(date), DATE_SUB(date, INTERVAL 1 YEAR) AS last_year
            FROM data_points CROSS JOIN (SELECT @sum := 0, @year := 0) AS init
            WHERE series_id = ? AND current = 1 ORDER BY date) AS t1
      LEFT JOIN (SELECT date, @sum := IF(@year = YEAR(date), @sum, 0) + value AS last_ytd,
                   @year := year(date), pseudo_history
                 FROM data_points CROSS JOIN (SELECT @sum := 0, @year := 0) AS init
                 WHERE series_id = ? AND current = 1 ORDER BY date) AS t2 ON (t1.last_year = t2.date)
      LEFT JOIN series ON t1.series_id = series.id;`,
		PlaceholderCount: 2,
		Label:            "ytd",
	},
	YTDPercentChange: { // ytd percent change from 1 year ago
		Statement: `SELECT t1.date, (t1.ytd/t2.last_ytd - 1)*100 AS ytd,
				(t1.pseudo_history = b'1') AND (t2.pseudo_history = b'1') AS ph
      FROM (SELECT date, value, @sum := IF(@year = YEAR(date), @sum, 0) + value AS ytd,
              @year := year(date), DATE_SUB(date, INTERVAL 1 YEAR) AS last_year, pseudo_history
            FROM data_points CROSS JOIN (SELECT @sum := 0, @year := 0) AS init
            WHERE series_id = ? AND current = 1 ORDER BY date) AS t1
      LEFT JOIN (SELECT date, @sum := IF(@year = YEAR(date), @sum, 0) + value AS last_ytd,
                   @year := year(date), pseudo_history
                 FROM data_points CROSS JOIN (SELECT @sum := 0, @year := 0) AS init
                 WHERE series_id = ? AND current = 1 ORDER BY date) AS t2 ON (t1.last_year = t2.date);`,
		PlaceholderCount: 2,
		Label:            "ytd",
	},
}

var seriesPrefix = `SELECT
	series.id, series.name, series.description, frequency, series.seasonally_adjusted, series.seasonal_adjustment,
	COALESCE(NULLIF(series.unitsLabel, ''), NULLIF(measurements.units_label, '')),
	COALESCE(NULLIF(series.unitsLabelShort, ''), NULLIF(measurements.units_label_short, '')),
	COALESCE(NULLIF(series.dataPortalName, ''), measurements.data_portal_name), measurements.percent, measurements.real,
	COALESCE(NULLIF(sources.description, ''), NULLIF(measurement_sources.description, '')),
	COALESCE(NULLIF(series.source_link, ''), NULLIF(measurements.source_link, ''), NULLIF(sources.link, ''), NULLIF(measurement_sources.link, '')),
	COALESCE(NULLIF(source_details.description, ''), NULLIF(measurement_source_details.description, '')),
	data_list_measurements.indent, series.base_year, series.decimals,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM series
	LEFT JOIN geographies ON series.name LIKE CONCAT('%@', geographies.handle, '.%')
	LEFT JOIN measurement_series ON measurement_series.series_id = series.id
	LEFT JOIN measurements ON measurements.id = measurement_series.measurement_id
	LEFT JOIN data_list_measurements ON data_list_measurements.measurement_id = measurements.id
	LEFT JOIN categories ON categories.data_list_id = data_list_measurements.data_list_id
	LEFT JOIN sources ON sources.id = series.source_id
	LEFT JOIN sources AS measurement_sources ON measurement_sources.id = measurements.source_id
	LEFT JOIN source_details ON source_details.id = series.source_detail_id
	LEFT JOIN source_details AS measurement_source_details ON measurement_source_details.id = measurements.source_detail_id
	LEFT JOIN feature_toggles ON feature_toggles.name = 'filter_by_quarantine'
	WHERE categories.id = ?
	AND NOT categories.hidden
	AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)`
var measurementSeriesPrefix = `SELECT DISTINCT
	series.id, series.name, series.description, frequency, series.seasonally_adjusted, series.seasonal_adjustment,
	COALESCE(NULLIF(series.unitsLabel, ''), NULLIF(measurements.units_label, '')),
	COALESCE(NULLIF(series.unitsLabelShort, ''), NULLIF(measurements.units_label_short, '')),
	COALESCE(NULLIF(series.dataPortalName, ''), measurements.data_portal_name), measurements.percent, measurements.real,
	COALESCE(NULLIF(sources.description, ''), NULLIF(measurement_sources.description, '')),
	COALESCE(NULLIF(series.source_link, ''), NULLIF(measurements.source_link, ''), NULLIF(sources.link, ''), NULLIF(measurement_sources.link, '')),
	COALESCE(NULLIF(source_details.description, ''), NULLIF(measurement_source_details.description, '')),
	NULL, series.base_year, series.decimals,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM measurements
	LEFT JOIN measurement_series ON measurement_series.measurement_id = measurements.id
	LEFT JOIN series ON series.id = measurement_series.series_id
	LEFT JOIN geographies ON series.name LIKE CONCAT('%@', geographies.handle, '.%')
	LEFT JOIN sources ON sources.id = series.source_id
	LEFT JOIN sources AS measurement_sources ON measurement_sources.id = measurements.source_id
	LEFT JOIN source_details ON source_details.id = series.source_detail_id
	LEFT JOIN source_details AS measurement_source_details ON measurement_source_details.id = measurements.source_detail_id
	LEFT JOIN feature_toggles ON feature_toggles.name = 'filter_by_quarantine'
	WHERE measurements.id = ? AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)`
var geoFilter = ` AND series.name LIKE CONCAT('%@', ? ,'.%') `
var freqFilter = ` AND series.name LIKE CONCAT('%@%.', ?) `
var sortStmt = ` ORDER BY data_list_measurements.list_order;`
var siblingsPrefix = `SELECT DISTINCT
        series.id, series.name, series.description, frequency, series.seasonally_adjusted, series.seasonal_adjustment,
	COALESCE(NULLIF(series.unitsLabel, ''), NULLIF(measurements.units_label, '')),
	COALESCE(NULLIF(series.unitsLabelShort, ''), NULLIF(measurements.units_label_short, '')),
	COALESCE(NULLIF(series.dataPortalName, ''), measurements.data_portal_name), measurements.percent, measurements.real,
	COALESCE(NULLIF(sources.description, ''), NULLIF(measurement_sources.description, '')),
	COALESCE(NULLIF(series.source_link, ''), NULLIF(measurements.source_link, ''), NULLIF(sources.link, ''), NULLIF(measurement_sources.link, '')),
	COALESCE(NULLIF(source_details.description, ''), NULLIF(measurement_source_details.description, '')),
	NULL, series.base_year, series.decimals,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM (SELECT measurement_id FROM measurement_series where series_id = ?) as measure
	LEFT JOIN measurements ON measurements.id = measure.measurement_id
	LEFT JOIN measurement_series ON measurement_series.measurement_id = measurements.id
	LEFT JOIN series ON series.id = measurement_series.series_id
	LEFT JOIN sources ON sources.id = series.source_id
	LEFT JOIN sources AS measurement_sources ON measurement_sources.id = measurements.source_id
	LEFT JOIN source_details ON source_details.id = series.source_detail_id
	LEFT JOIN source_details AS measurement_source_details ON measurement_source_details.id = measurements.source_detail_id
	LEFT JOIN geographies ON series.name LIKE CONCAT('%@', geographies.handle, '.%')
	LEFT JOIN feature_toggles ON feature_toggles.name = 'filter_by_quarantine'
	WHERE NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)`

func (r *SeriesRepository) GetSeriesByGroupAndFreq(
	groupId int64,
	freq string,
	groupType GroupType,
) (seriesList []models.DataPortalSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = ";"
	}
	rows, err := r.DB.Query(
		strings.Join([]string{prefix, freqFilter, sort}, ""),
		groupId,
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

func (r *SeriesRepository) GetSeriesByGroupGeoAndFreq(
	groupId int64,
	geoHandle string,
	freq string,
	groupType GroupType,
) (seriesList []models.DataPortalSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = ";"
	}
	rows, err := r.DB.Query(
		strings.Join([]string{prefix, geoFilter, freqFilter, sort}, ""),
		groupId,
		geoHandle,
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

func (r *SeriesRepository) GetInflatedSeriesByGroupGeoAndFreq(
	groupId int64,
	geoHandle string,
	freq string,
	groupType GroupType,
) (seriesList []models.InflatedSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = ";"
	}
	rows, err := r.DB.Query(
		strings.Join([]string{prefix, geoFilter, freqFilter, sort}, ""),
		groupId,
		geoHandle,
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

func (r *SeriesRepository) GetSeriesByGroupAndGeo(
	groupId int64,
	geoHandle string,
	groupType GroupType,
) (seriesList []models.DataPortalSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = ";"
	}
	rows, err := r.DB.Query(
		strings.Join([]string{prefix, geoFilter, sort}, ""),
		groupId,
		geoHandle,
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

func (r *SeriesRepository) GetInflatedSeriesByGroup(
	groupId int64,
	groupType GroupType,
) (seriesList []models.InflatedSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = ";"
	}
	rows, err := r.DB.Query(
		strings.Join([]string{prefix, sort}, ""),
		groupId,
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

func (r *SeriesRepository) GetSeriesByGroup(
	groupId int64,
	groupType GroupType,
) (seriesList []models.DataPortalSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = ";"
	}
	rows, err := r.DB.Query(
		strings.Join([]string{prefix, sort}, ""),
		groupId,
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

func (r *SeriesRepository) GetFreqByCategory(categoryId int64) (frequencies []models.FrequencyResult, err error) {
	rows, err := r.DB.Query(`SELECT DISTINCT(RIGHT(series.name, 1)) as freq
	FROM categories
	LEFT JOIN data_list_measurements ON data_list_measurements.data_list_id = categories.data_list_id
	LEFT JOIN measurement_series ON measurement_series.measurement_id = data_list_measurements.measurement_id
	LEFT JOIN series ON series.id = measurement_series.series_id
	LEFT JOIN feature_toggles ON feature_toggles.name = 'filter_by_quarantine'
	WHERE categories.id = ?
	AND NOT categories.hidden
	AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)
	ORDER BY FIELD(freq, "A", "S", "Q", "M", "W", "D");`, categoryId)
	if err != nil {
		return
	}
	for rows.Next() {
		frequency := models.Frequency{}
		err = rows.Scan(
			&frequency.Freq,
		)
		if err != nil {
			return
		}
		frequencies = append(
			frequencies,
			models.FrequencyResult{Freq: frequency.Freq, Label: freqLabel[frequency.Freq]},
		)
	}
	return

}

func (r *SeriesRepository) GetSeriesSiblingsById(seriesId int64) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(siblingsPrefix, seriesId)
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

func (r *SeriesRepository) GetSeriesSiblingsByIdAndFreq(
	seriesId int64,
	freq string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(strings.Join([]string{siblingsPrefix, freqFilter}, ""), seriesId, freq)
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

func (r *SeriesRepository) GetSeriesSiblingsByIdAndGeo(
	seriesId int64,
	geo string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(strings.Join([]string{siblingsPrefix, geoFilter}, ""), seriesId, geo)
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

func (r *SeriesRepository) GetSeriesSiblingsByIdGeoAndFreq(
	seriesId int64,
	geo string,
	freq string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(
		strings.Join([]string{siblingsPrefix, geoFilter, freqFilter}, ""),
		seriesId, geo, freq)
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

func (r *SeriesRepository) GetSeriesSiblingsFreqById(
	seriesId int64,
) (frequencyList []models.FrequencyResult, err error) {
	rows, err := r.DB.Query(`SELECT DISTINCT(RIGHT(series.name, 1)) as freq
	FROM series
	JOIN (SELECT name FROM series where id = ?) as original_series
	LEFT JOIN feature_toggles ON feature_toggles.name = 'filter_by_quarantine'
	WHERE series.name LIKE CONCAT(TRIM(TRAILING 'NS' FROM LEFT(original_series.name, LOCATE("@", original_series.name))), '%')
	AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)
	ORDER BY FIELD(freq, "A", "S", "Q", "M", "W", "D");`, seriesId)
	if err != nil {
		return
	}
	for rows.Next() {
		frequency := models.Frequency{}
		err = rows.Scan(
			&frequency.Freq,
		)
		if err != nil {
			return
		}
		frequencyList = append(
			frequencyList,
			models.FrequencyResult{Freq: frequency.Freq, Label: freqLabel[frequency.Freq]},
		)
	}
	return
}

func (r *SeriesRepository) GetSeriesById(seriesId int64) (dataPortalSeries models.DataPortalSeries, err error) {
	row := r.DB.QueryRow(`SELECT DISTINCT
	series.id, series.name, series.description, frequency, series.seasonally_adjusted, series.seasonal_adjustment,
	COALESCE(NULLIF(series.unitsLabel, ''), NULLIF(measurements.units_label, '')),
	COALESCE(NULLIF(series.unitsLabelShort, ''), NULLIF(measurements.units_label_short, '')),
	COALESCE(NULLIF(series.dataPortalName, ''), measurements.data_portal_name), measurements.percent, measurements.real,
	COALESCE(NULLIF(sources.description, ''), NULLIF(measurement_sources.description, '')),
	COALESCE(NULLIF(series.source_link, ''), NULLIF(measurements.source_link, ''), NULLIF(sources.link, ''), NULLIF(measurement_sources.link, '')),
	COALESCE(NULLIF(source_details.description, ''), NULLIF(measurement_source_details.description, '')),
	series.base_year, series.decimals,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM series
	LEFT JOIN geographies ON series.name LIKE CONCAT('%@', geographies.handle, '.%')
	LEFT JOIN measurement_series ON measurement_series.series_id = series.id
	LEFT JOIN measurements ON measurements.id = measurement_series.measurement_id
	LEFT JOIN sources ON sources.id = series.source_id
	LEFT JOIN sources AS measurement_sources ON measurement_sources.id = measurements.source_id
	LEFT JOIN source_details ON source_details.id = series.source_detail_id
	LEFT JOIN source_details AS measurement_source_details ON measurement_source_details.id = measurements.source_detail_id
	LEFT JOIN feature_toggles ON feature_toggles.name = 'filter_by_quarantine'
	WHERE series.id = ? AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined);`, seriesId)
	dataPortalSeries, err = getNextSeriesFromRow(row)
	if err != nil {
		return
	}
	geoFreqs, freqGeos, err := getFreqGeoCombinations(r, seriesId)
	if err != nil {
		return
	}
	dataPortalSeries.GeographyFrequencies = &geoFreqs
	dataPortalSeries.FrequencyGeographies = &freqGeos
	dataPortalSeries.GeoFreqsDeprecated = &geoFreqs
	dataPortalSeries.FreqGeosDeprecated = &freqGeos
	return
}

// GetSeriesObservations returns an observations struct containing the default transformations.
// It checks the value of percent for the selected series and chooses the appropriate transformations.
func (r *SeriesRepository) GetSeriesObservations(
	seriesId int64,
) (seriesObservations models.SeriesObservations, err error) {
	var start, end time.Time
	var percent sql.NullBool
	YOY, YTD := YOYPercentChange, YTDPercentChange

	err = r.DB.QueryRow(`SELECT measurements.percent
	FROM series
	LEFT JOIN measurement_series ON measurement_series.series_id = series.id
	LEFT JOIN measurements ON measurements.id = measurement_series.measurement_id
	LEFT JOIN feature_toggles ON feature_toggles.name = 'filter_by_quarantine'
	WHERE series.id = ? AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)`, seriesId).Scan(&percent)
	if err != nil {
		return
	}
	if percent.Valid && percent.Bool {
		YOY = YOYChange
		YTD = YTDChange
	}

	lvlTransform, err := r.GetTransformation(Levels, seriesId, &start, &end)
	if err != nil {
		return
	}
	yoyTransform, err := r.GetTransformation(YOY, seriesId, &start, &end)
	if err != nil {
		return
	}
	ytdTransform, err := r.GetTransformation(YTD, seriesId, &start, &end)
	if err != nil {
		return
	}

	seriesObservations.TransformationResults = []models.TransformationResult{lvlTransform, yoyTransform, ytdTransform}
	seriesObservations.ObservationStart = start
	seriesObservations.ObservationEnd = end
	return
}

func variadicSeriesId(seriesId int64, count int) []interface{} {
	variadic := make([]interface{}, count, count)
	for i := range variadic {
		variadic[i] = seriesId
	}
	return variadic
}

func (r *SeriesRepository) GetTransformation(
	transformation string,
	seriesId int64,
	currentStart *time.Time,
	currentEnd *time.Time,
) (
	transformationResult models.TransformationResult,
	err error,
) {
	var observationStart, observationEnd time.Time
	rows, err := r.DB.Query(
		transformations[transformation].Statement,
		variadicSeriesId(seriesId, transformations[transformation].PlaceholderCount)...,
	)
	if err != nil {
		return
	}
	var (
		observations []models.DataPortalObservation
	)

	for rows.Next() {
		observation := models.Observation{}
		err = rows.Scan(
			&observation.Date,
			&observation.Value,
			&observation.PseudoHistory,
		)
		if err != nil {
			return
		}
		if !observation.Value.Valid {
			continue
		}
		if observationStart.IsZero() || observation.Date.Before(observationStart) {
			observationStart = observation.Date
		}
		if observationEnd.IsZero() || observationEnd.Before(observation.Date) {
			observationEnd = observation.Date
		}
		dataPortalObservation := models.DataPortalObservation{
			Date:  observation.Date,
			Value: observation.Value.Float64,
		}
		if observation.PseudoHistory.Valid && observation.PseudoHistory.Bool {
			dataPortalObservation.PseudoHistory = &observation.PseudoHistory.Bool
		}
		observations = append(
			observations,
			dataPortalObservation,
		)
	}
	if currentStart.IsZero() || (!observationStart.IsZero() && currentStart.After(observationStart)) {
		*currentStart = observationStart
	}
	if currentEnd.IsZero() || (!observationStart.IsZero() && currentEnd.Before(observationEnd)) {
		*currentEnd = observationEnd
	}
	transformationResult.Transformation = transformations[transformation].Label
	transformationResult.Observations = observations
	return
}
