package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"strconv"
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
	Levels            = "lvl"
	YOYPercentChange  = "pc1"
	YTDPercentChange  = "ytdpc1"
	C5MAPercentChange = "c5mapc1"
	YOYChange         = "ch1"
	YTDChange         = "ytdch1"
	C5MAChange        = "c5mach1"
)

var transformations map[string]transformation = map[string]transformation{
	Levels: { // untransformed value
		Statement: `SELECT date, value/units, (pseudo_history = b'1'), series.decimals
		FROM public_data_points
		LEFT JOIN series ON public_data_points.series_id = series.id
		WHERE series_id = ?;`,
		PlaceholderCount: 1,
		Label:            "lvl",
	},
	YOYPercentChange: { // percent change from 1 year ago
		Statement: `SELECT t1.date, (t1.value/t2.last_value - 1)*100 AS yoy,
				(t1.pseudo_history = b'1') AND (t2.pseudo_history = b'1') AS ph, series.decimals
				FROM (SELECT series_id, value, date, pseudo_history, DATE_SUB(date, INTERVAL 1 YEAR) AS last_year
				      FROM public_data_points WHERE series_id = ?) AS t1
				LEFT JOIN (SELECT value AS last_value, date, pseudo_history
				           FROM public_data_points WHERE series_id = ?) AS t2 ON (t1.last_year = t2.date)
				LEFT JOIN series ON t1.series_id = series.id;`,
		PlaceholderCount: 2,
		Label:            "pc1",
	},
	YOYChange: { // change from 1 year ago
		Statement: `SELECT t1.date, (t1.value - t2.last_value)/series.units AS yoy,
				(t1.pseudo_history = b'1') AND (t2.pseudo_history = b'1') AS ph, series.decimals
				FROM (SELECT series_id, value, date, pseudo_history, DATE_SUB(date, INTERVAL 1 YEAR) AS last_year
				      FROM public_data_points WHERE series_id = ?) AS t1
				LEFT JOIN (SELECT value AS last_value, date, pseudo_history
				           FROM public_data_points WHERE series_id = ?) AS t2 ON (t1.last_year = t2.date)
				LEFT JOIN series ON t1.series_id = series.id;`,
		PlaceholderCount: 2,
		Label:            "pc1",
	},
	YTDChange: { // ytd change from 1 year ago
		Statement: `SELECT t1.date, (t1.ytd/t1.count - t2.last_ytd/t2.last_count)/series.units AS ytd,
				(t1.pseudo_history = b'1') AND (t2.pseudo_history = b'1') AS ph, series.decimals
	  FROM (SELECT date, value, series_id, pseudo_history, @sum := IF(@year = YEAR(date), @sum, 0) + value AS ytd,
	  		  @count := IF(@year = YEAR(date), @count, 0) + 1 AS count,
              @year := year(date), DATE_SUB(date, INTERVAL 1 YEAR) AS last_year
            FROM public_data_points CROSS JOIN (SELECT @sum := 0, @year := 0) AS init
            WHERE series_id = ? ORDER BY date) AS t1
	  LEFT JOIN (SELECT date, @sum := IF(@year = YEAR(date), @sum, 0) + value AS last_ytd,
				   @count := IF(@year = YEAR(date), @count, 0) + 1 AS last_count,
                   @year := year(date), pseudo_history
                 FROM public_data_points CROSS JOIN (SELECT @sum := 0, @year := 0) AS init
                 WHERE series_id = ? ORDER BY date) AS t2 ON (t1.last_year = t2.date)
      LEFT JOIN series ON t1.series_id = series.id;`,
		PlaceholderCount: 2,
		Label:            "ytd",
	},
	YTDPercentChange: { // ytd percent change from 1 year ago
		Statement: `SELECT t1.date, (t1.ytd/t2.last_ytd - 1)*100 AS ytd,
				(t1.pseudo_history = b'1') AND (t2.pseudo_history = b'1') AS ph, series.decimals
      FROM (SELECT series_id, date, value, @sum := IF(@year = YEAR(date), @sum, 0) + value AS ytd,
              @year := year(date), DATE_SUB(date, INTERVAL 1 YEAR) AS last_year, pseudo_history
            FROM public_data_points CROSS JOIN (SELECT @sum := 0, @year := 0) AS init
            WHERE series_id = ? ORDER BY date) AS t1
      LEFT JOIN (SELECT date, @sum := IF(@year = YEAR(date), @sum, 0) + value AS last_ytd,
                   @year := year(date), pseudo_history
                 FROM public_data_points CROSS JOIN (SELECT @sum := 0, @year := 0) AS init
                 WHERE series_id = ? ORDER BY date) AS t2 ON (t1.last_year = t2.date)
      LEFT JOIN series ON t1.series_id = series.id;`,
		PlaceholderCount: 2,
		Label:            "ytd",
	},
	C5MAPercentChange: { // c5ma percent change from 1 year ago
		Statement: `SELECT t1.date, (t1.c5ma/t2.last_c5ma - 1)*100 AS c5ma, 
			(t1.pseudo_history = b'1') AND (t2.pseudo_history = b'1') AS ph, series.decimals
			FROM (SELECT pdp2.series_id, pdp1.date, CASE WHEN count(*) = 5 THEN avg(pdp2.value) ELSE NULL END AS c5ma, DATE_SUB(pdp1.date, INTERVAL 1 YEAR) AS last_year, pdp1.pseudo_history FROM
				(SELECT date, value, pseudo_history FROM public_data_points WHERE series_id = ?) AS pdp1
				INNER JOIN public_data_points AS pdp2 ON pdp2.date BETWEEN DATE_SUB(pdp1.date, INTERVAL 2 YEAR) AND DATE_ADD(pdp1.date, INTERVAL 2 YEAR) WHERE series_id = ?
				GROUP by series_id, date, last_year, pseudo_history) AS t1
			LEFT JOIN (SELECT pdp2.series_id, pdp1.date, CASE WHEN count(*) = 5 THEN avg(pdp2.value) ELSE NULL END AS last_c5ma, pdp1.pseudo_history FROM
				(SELECT date, value, pseudo_history FROM public_data_points WHERE series_id = ?) AS pdp1
				INNER JOIN public_data_points AS pdp2 ON pdp2.date BETWEEN DATE_SUB(pdp1.date, INTERVAL 2 YEAR) AND DATE_ADD(pdp1.date, INTERVAL 2 YEAR) WHERE series_id = ?
				GROUP by series_id, date, pseudo_history) AS t2 ON (t1.last_year = t2.date)
      			LEFT JOIN series ON t1.series_id = series.id;`,
		PlaceholderCount: 4,
		Label:            "c5ma",
	},
	C5MAChange: { // cm5a change from 1 year ago
		Statement: `SELECT t1.date, (t1.c5ma - t2.last_c5ma)/series.units AS c5ma, 
			(t1.pseudo_history = b'1') AND (t2.pseudo_history = b'1') AS ph, series.decimals
			FROM (SELECT pdp2.series_id, pdp1.date, CASE WHEN count(*) = 5 THEN avg(pdp2.value) ELSE NULL END AS c5ma, DATE_SUB(pdp1.date, INTERVAL 1 YEAR) AS last_year, pdp1.pseudo_history FROM
				(SELECT date, value, pseudo_history FROM public_data_points WHERE series_id = ?) AS pdp1
				INNER JOIN public_data_points AS pdp2 ON pdp2.date BETWEEN DATE_SUB(pdp1.date, INTERVAL 2 YEAR) AND DATE_ADD(pdp1.date, INTERVAL 2 YEAR) WHERE series_id = ?
				GROUP by series_id, date, last_year, pseudo_history) AS t1
			LEFT JOIN (SELECT pdp2.series_id, pdp1.date, CASE WHEN count(*) = 5 THEN avg(pdp2.value) ELSE NULL END AS last_c5ma, pdp1.pseudo_history FROM
				(SELECT date, value, pseudo_history FROM public_data_points WHERE series_id = ?) AS pdp1
				INNER JOIN public_data_points AS pdp2 ON pdp2.date BETWEEN DATE_SUB(pdp1.date, INTERVAL 2 YEAR) AND DATE_ADD(pdp1.date, INTERVAL 2 YEAR) WHERE series_id = ?
				GROUP by series_id, date, pseudo_history) AS t2 ON (t1.last_year = t2.date)
			LEFT JOIN series ON t1.series_id = series.id;`,
		PlaceholderCount: 4,
		Label:            "c5ma",
	},
}

var seriesPrefix = `SELECT
	series.id, series.name, series.universe, series.description, frequency, series.seasonally_adjusted, series.seasonal_adjustment,
	COALESCE(NULLIF(units.long_label, ''), NULLIF(MAX(measurement_units.long_label), '')),
	COALESCE(NULLIF(units.short_label, ''), NULLIF(MAX(measurement_units.short_label), '')),
	COALESCE(NULLIF(series.dataPortalName, ''), MAX(measurements.data_portal_name)), series.percent, series.real,
	COALESCE(NULLIF(sources.description, ''), NULLIF(MAX(measurement_sources.description), '')),
	COALESCE(NULLIF(series.source_link, ''), NULLIF(MAX(measurements.source_link), ''), NULLIF(sources.link, ''), NULLIF(MAX(measurement_sources.link), '')),
	COALESCE(NULLIF(source_details.description, ''), NULLIF(MAX(measurement_source_details.description), '')),
	MAX(measurements.table_prefix), MAX(measurements.table_postfix),
	MAX(measurements.id), MAX(measurements.data_portal_name),
	MAX(data_list_measurements.indent), series.base_year, series.decimals,
	MAX(geographies.fips), MAX(geographies.handle) AS shandle, MAX(geographies.display_name), MAX(geographies.display_name_short)
	FROM series
	LEFT JOIN measurement_series ON measurement_series.series_id = series.id
	LEFT JOIN measurements ON measurements.id = measurement_series.measurement_id
	LEFT JOIN data_list_measurements ON data_list_measurements.measurement_id = measurements.id
	LEFT JOIN categories ON categories.data_list_id = data_list_measurements.data_list_id
	LEFT JOIN geographies ON geographies.id = series.geography_id
	LEFT JOIN units ON units.id = series.unit_id
	LEFT JOIN units AS measurement_units ON measurement_units.id = measurements.unit_id
	LEFT JOIN sources ON sources.id = series.source_id
	LEFT JOIN sources AS measurement_sources ON measurement_sources.id = measurements.source_id
	LEFT JOIN source_details ON source_details.id = series.source_detail_id
	LEFT JOIN source_details AS measurement_source_details ON measurement_source_details.id = measurements.source_detail_id
	LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE categories.id = ?
	AND NOT (categories.hidden OR categories.masked)
	AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)`
var measurementSeriesPrefix = `SELECT
	series.id, series.name, series.universe, series.description, frequency, series.seasonally_adjusted, series.seasonal_adjustment,
	COALESCE(NULLIF(units.long_label, ''), NULLIF(MAX(measurement_units.long_label), '')),
	COALESCE(NULLIF(units.short_label, ''), NULLIF(MAX(measurement_units.short_label), '')),
	COALESCE(NULLIF(series.dataPortalName, ''), MAX(measurements.data_portal_name)), series.percent, series.real,
	COALESCE(NULLIF(sources.description, ''), NULLIF(MAX(measurement_sources.description), '')),
	COALESCE(NULLIF(series.source_link, ''), NULLIF(MAX(measurements.source_link), ''), NULLIF(sources.link, ''), NULLIF(MAX(measurement_sources.link), '')),
	COALESCE(NULLIF(source_details.description, ''), NULLIF(MAX(measurement_source_details.description), '')),
	MAX(measurements.table_prefix), MAX(measurements.table_postfix),
	MAX(measurements.id), MAX(measurements.data_portal_name),
	NULL, series.base_year, series.decimals,
	MAX(geographies.fips), MAX(geographies.handle) AS shandle, MAX(geographies.display_name), MAX(geographies.display_name_short)
	FROM measurements
	LEFT JOIN measurement_series ON measurement_series.measurement_id = measurements.id
	LEFT JOIN series ON series.id = measurement_series.series_id
	LEFT JOIN geographies ON geographies.id = series.geography_id
	LEFT JOIN units ON units.id = series.unit_id
	LEFT JOIN units AS measurement_units ON measurement_units.id = measurements.unit_id
	LEFT JOIN sources ON sources.id = series.source_id
	LEFT JOIN sources AS measurement_sources ON measurement_sources.id = measurements.source_id
	LEFT JOIN source_details ON source_details.id = series.source_detail_id
	LEFT JOIN source_details AS measurement_source_details ON measurement_source_details.id = measurements.source_detail_id
	LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE measurements.id = ? AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)`
var geoFilter = ` AND geographies.handle = UPPER(?) `
var freqFilter = ` AND series.frequency = ? `
var measurementPostfix = ` GROUP BY series.id;`
var sortStmt = ` GROUP BY series.id ORDER BY MAX(data_list_measurements.list_order);`
var siblingSortStmt = ` GROUP BY series.id;`
var siblingsPrefix = `SELECT
    series.id, series.name, series.universe, series.description, frequency, series.seasonally_adjusted, series.seasonal_adjustment,
	COALESCE(NULLIF(units.long_label, ''), NULLIF(MAX(measurement_units.long_label), '')),
	COALESCE(NULLIF(units.short_label, ''), NULLIF(MAX(measurement_units.short_label), '')),
	COALESCE(NULLIF(series.dataPortalName, ''), MAX(measurements.data_portal_name)), series.percent, series.real,
	COALESCE(NULLIF(sources.description, ''), NULLIF(MAX(measurement_sources.description), '')),
	COALESCE(NULLIF(series.source_link, ''), NULLIF(MAX(measurements.source_link), ''), NULLIF(sources.link, ''), NULLIF(MAX(measurement_sources.link), '')),
	COALESCE(NULLIF(source_details.description, ''), NULLIF(MAX(measurement_source_details.description), '')),
	MAX(measurements.table_prefix), MAX(measurements.table_postfix),
	MAX(measurements.id), MAX(measurements.data_portal_name),
	NULL, series.base_year, series.decimals,
	MAX(geographies.fips), MAX(geographies.handle) AS shandle, MAX(geographies.display_name), MAX(geographies.display_name_short)
	FROM (SELECT measurement_id FROM measurement_series where series_id = ?) as measure
	LEFT JOIN measurements ON measurements.id = measure.measurement_id
	LEFT JOIN measurement_series ON measurement_series.measurement_id = measurements.id
	LEFT JOIN series ON series.id = measurement_series.series_id
	LEFT JOIN public_data_points ON public_data_points.series_id = series.id
	LEFT JOIN geographies ON geographies.id = series.geography_id
	LEFT JOIN units ON units.id = series.unit_id
	LEFT JOIN units AS measurement_units ON measurement_units.id = measurements.unit_id
	LEFT JOIN sources ON sources.id = series.source_id
	LEFT JOIN sources AS measurement_sources ON measurement_sources.id = measurements.source_id
	LEFT JOIN source_details ON source_details.id = series.source_detail_id
	LEFT JOIN source_details AS measurement_source_details ON measurement_source_details.id = measurements.source_detail_id
	LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE NOT series.restricted
	AND public_data_points.value IS NOT NULL
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
		sort = measurementPostfix
	}
	rows, err := r.DB.Query(
		strings.Join([]string{prefix, freqFilter, sort}, ""),
		groupId,
		freqDbNames[strings.ToUpper(freq)],
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
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
		sort = measurementPostfix
	}
	rows, err := r.DB.Query(
		strings.Join([]string{prefix, geoFilter, freqFilter, sort}, ""),
		groupId,
		geoHandle,
		freqDbNames[strings.ToUpper(freq)],
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
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
		sort = measurementPostfix
	}
	rows, err := r.DB.Query(
		strings.Join([]string{prefix, geoFilter, freqFilter, sort}, ""),
		groupId,
		geoHandle,
		freqDbNames[strings.ToUpper(freq)],
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
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
		sort = measurementPostfix
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
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
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
		sort = measurementPostfix
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
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
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
		sort = measurementPostfix
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
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetFreqByCategory(categoryId int64) (frequencies []models.DataPortalFrequency, err error) {
	rows, err := r.DB.Query(`SELECT DISTINCT(RIGHT(series.name, 1)) as freq
	FROM categories
	LEFT JOIN data_list_measurements ON data_list_measurements.data_list_id = categories.data_list_id
	LEFT JOIN measurement_series ON measurement_series.measurement_id = data_list_measurements.measurement_id
	LEFT JOIN series ON series.id = measurement_series.series_id
	LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE categories.id = ?
	AND NOT (categories.hidden OR categories.masked)
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
			models.DataPortalFrequency{Freq: frequency.Freq, Label: freqLabel[frequency.Freq]},
		)
	}
	return

}

func (r *SeriesRepository) GetSeriesSiblingsById(seriesId int64) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(
		strings.Join([]string{siblingsPrefix, siblingSortStmt}, ""),
		seriesId,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesSiblingsByIdAndFreq(
	seriesId int64,
	freq string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(strings.Join(
		[]string{siblingsPrefix, freqFilter, siblingSortStmt}, ""),
		seriesId,
		freqDbNames[strings.ToUpper(freq)],
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesSiblingsByIdAndGeo(
	seriesId int64,
	geo string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(
		strings.Join([]string{siblingsPrefix, geoFilter, siblingSortStmt}, ""),
		seriesId,
		geo,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
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
		strings.Join([]string{siblingsPrefix, geoFilter, freqFilter, siblingSortStmt}, ""),
		seriesId,
		geo,
		freqDbNames[strings.ToUpper(freq)],
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesSiblingsFreqById(
	seriesId int64,
) (frequencyList []models.DataPortalFrequency, err error) {
	rows, err := r.DB.Query(`SELECT DISTINCT(RIGHT(series.name, 1)) as freq
	FROM series
	JOIN (SELECT name, universe FROM series WHERE id = ?) as original_series
	LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE series.universe = original_series.universe
	AND TRIM(TRAILING 'NS' FROM TRIM(TRAILING '&' FROM SUBSTRING_INDEX(series.name, '@', 1))) =
	    TRIM(TRAILING 'NS' FROM TRIM(TRAILING '&' FROM SUBSTRING_INDEX(original_series.name, '@', 1)))
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
			models.DataPortalFrequency{Freq: frequency.Freq, Label: freqLabel[frequency.Freq]},
		)
	}
	return
}

func (r *SeriesRepository) GetSeriesById(seriesId int64) (dataPortalSeries models.DataPortalSeries, err error) {
	row, err := r.DB.Query(`SELECT DISTINCT
	series.id, series.name, series.universe, series.description, frequency, series.seasonally_adjusted, series.seasonal_adjustment,
	COALESCE(NULLIF(units.long_label, ''), NULLIF(measurement_units.long_label, '')),
	COALESCE(NULLIF(units.short_label, ''), NULLIF(measurement_units.short_label, '')),
	COALESCE(NULLIF(series.dataPortalName, ''), measurements.data_portal_name), series.percent, series.real,
	COALESCE(NULLIF(sources.description, ''), NULLIF(measurement_sources.description, '')),
	COALESCE(NULLIF(series.source_link, ''), NULLIF(measurements.source_link, ''), NULLIF(sources.link, ''), NULLIF(measurement_sources.link, '')),
	COALESCE(NULLIF(source_details.description, ''), NULLIF(measurement_source_details.description, '')),
	measurements.table_prefix, measurements.table_postfix,
	measurements.id, measurements.data_portal_name,
	NULL, series.base_year, series.decimals,
	geo.fips, geo.handle AS shandle, geo.display_name, geo.display_name_short
	FROM series
	LEFT JOIN geographies geo ON geo.id = series.geography_id
	LEFT JOIN measurement_series ON measurement_series.series_id = series.id
	LEFT JOIN measurements ON measurements.id = measurement_series.measurement_id
	LEFT JOIN units ON units.id = series.unit_id
	LEFT JOIN units AS measurement_units ON measurement_units.id = measurements.unit_id
	LEFT JOIN sources ON sources.id = series.source_id
	LEFT JOIN sources AS measurement_sources ON measurement_sources.id = measurements.source_id
	LEFT JOIN source_details ON source_details.id = series.source_detail_id
	LEFT JOIN source_details AS measurement_source_details ON measurement_source_details.id = measurements.source_detail_id
	LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE series.id = ? AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined);`, seriesId)
	if err != nil {
		return
	}
	for row.Next() {
		dataPortalSeries, err = getNextSeriesFromRows(row)
		if err != nil {
			return
		}
		geos, freqs, err := getAllFreqsGeos(r, seriesId)
		if err != nil {
			return dataPortalSeries, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		break
	}
	return
}

// GetSeriesObservations returns an observations struct containing the default transformations.
// It checks the value of percent for the selected series and chooses the appropriate transformations.
func (r *SeriesRepository) GetSeriesObservations(
	seriesId int64,
) (seriesObservations models.SeriesObservations, err error) {
	var start, end time.Time
	var percent sql.NullBool
	var universe string
	YOY, YTD, C5MA := YOYPercentChange, YTDPercentChange, C5MAPercentChange

	err = r.DB.QueryRow(`SELECT series.universe, series.percent
	FROM series
	LEFT JOIN feature_toggles ON feature_toggles.universe = series.universe AND feature_toggles.name = 'filter_by_quarantine'
	WHERE series.id = ? AND NOT series.restricted
	AND (feature_toggles.status IS NULL OR NOT feature_toggles.status OR NOT series.quarantined)`, seriesId).Scan(&universe, &percent)
	if err != nil {
		return
	}
	if percent.Valid && percent.Bool {
		YOY = YOYChange
		YTD = YTDChange
		C5MA = C5MAChange
	}

	transform, err := r.GetTransformation(Levels, seriesId, &start, &end)
	if err != nil {
		return
	}
	seriesObservations.TransformationResults = append(seriesObservations.TransformationResults, transform)

	if universe != "NTA" {
		transform, err = r.GetTransformation(YOY, seriesId, &start, &end)
		if err != nil {
			return
		}
		seriesObservations.TransformationResults = append(seriesObservations.TransformationResults, transform)

		transform, err = r.GetTransformation(YTD, seriesId, &start, &end)
		if err != nil {
			return
		}
		seriesObservations.TransformationResults = append(seriesObservations.TransformationResults, transform)
	}
	if universe == "NTA" {
		transform, err = r.GetTransformation(C5MA, seriesId, &start, &end)
		if err != nil {
			return
		}
		seriesObservations.TransformationResults = append(seriesObservations.TransformationResults, transform)
	}
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
		obsDates      []string
		obsValues     []string
		obsPseudoHist []bool
	)

	for rows.Next() {
		observation := models.Observation{}
		err = rows.Scan(
			&observation.Date,
			&observation.Value,
			&observation.PseudoHistory,
			&observation.Decimals,
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
		obsDates = append(obsDates, observation.Date.Format("2006-01-02"))
		obsValues = append(obsValues, strconv.FormatFloat(observation.Value.Float64, 'f', observation.Decimals, 64))
		obsPseudoHist = append(obsPseudoHist, observation.PseudoHistory.Bool)
	}
	if currentStart.IsZero() || (!observationStart.IsZero() && currentStart.After(observationStart)) {
		*currentStart = observationStart
	}
	if currentEnd.IsZero() || (!observationStart.IsZero() && currentEnd.Before(observationEnd)) {
		*currentEnd = observationEnd
	}
	transformationResult.Transformation = transformations[transformation].Label
	transformationResult.ObservationDates = obsDates
	transformationResult.ObservationValues = obsValues
	transformationResult.ObservationPHist = obsPseudoHist
	return
}

func (r *SeriesRepository) CreateSeriesPackage(
	id int64,
	universe string,
	categoryRepository *CategoryRepository,
)  (pkg models.DataPortalSeriesPackage, err error) {

	series, err := r.GetSeriesById(id)
	if err != nil {
		return
	}
	pkg.Series = series

	categories, err := categoryRepository.GetAllCategoriesByUniverse(universe)
	if err != nil {
		return
	}
	pkg.Categories = categories

	observations, err := r.GetSeriesObservations(id)
	if err != nil {
		return
	}
	pkg.Observations = observations

	siblings, err := r.GetSeriesSiblingsById(id)
	if err != nil {
		return
	}
	pkg.Siblings = siblings
	return
}

func (r *SeriesRepository) CreateAnalyzerPackage(
	ids []int64,
	universe string,
	categoryRepository *CategoryRepository,
) (pkg models.DataPortalAnalyzerPackage, err error) {

	pkg.InflatedSeries = []models.InflatedSeries{}

	for _, id := range ids {
		series, anErr := r.GetSeriesById(id)
		if anErr != nil {
			err = anErr
			return
		}
		observations, anErr := r.GetSeriesObservations(id)
		if anErr != nil {
			err = anErr
			return
		}
		pkg.InflatedSeries = append(pkg.InflatedSeries, models.InflatedSeries{DataPortalSeries: series, Observations: observations})
	}

	categories, err := categoryRepository.GetAllCategoriesByUniverse(universe)
	if err != nil {
		return
	}
	pkg.Categories = categories
	return
}
