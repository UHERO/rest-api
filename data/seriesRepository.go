package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"time"
	"strings"
)

type SeriesRepository struct {
	DB *sql.DB
}

type transformation struct {
	Statement string
	PlaceholderCount int
}

const (
	Levels = "lvl"
	YoyPCh = "pc1"
	YTD = "ytd"
)

var transformations map[string]transformation = map[string]transformation{
	"lvl": {
		Statement: `SELECT date, value FROM data_points WHERE series_id = ? and current = 1;`,
		PlaceholderCount: 1,
	},
	"pc1": {
		Statement: `SELECT t1.date, (t1.value/t2.last_value - 1)*100 AS yoy
				FROM (SELECT value, date, DATE_SUB(date, INTERVAL 1 YEAR) AS last_year
				FROM data_points WHERE series_id = ? AND current = 1) AS t1
				LEFT JOIN (SELECT value AS last_value, date
				FROM data_points WHERE series_id = ? and current = 1) AS t2
				ON (t1.last_year = t2.date);`,
		PlaceholderCount: 2,
	},
	"ytd": {
		Statement: `SELECT t1.date, (t1.ytd/t2.last_ytd - 1)*100 AS ytd
      FROM (SELECT date, value, @sum := IF(@year = YEAR(date), @sum, 0) + value AS ytd,
            @year := year(date), DATE_SUB(date, INTERVAL 1 YEAR) AS last_year
          FROM data_points CROSS JOIN (SELECT @sum := 0, @year := 0) AS init
          WHERE series_id = ? AND current = 1 ORDER BY date) AS t1
      LEFT JOIN (SELECT date, @sum := IF(@year = YEAR(date), @sum, 0) + value AS last_ytd,
            @year := year(date)
          FROM data_points CROSS JOIN (SELECT @sum := 0, @year := 0) AS init
          WHERE series_id = ? AND current = 1 ORDER BY date) AS t2
      ON (t1.last_year = t2.date);`,
		PlaceholderCount: 2,
	},
}

var seriesPrefix = `SELECT series.id, name, description, frequency,
	!(name REGEXP '.*NS@.*') AS seasonally_adjusted,
	unitsLabel, unitsLabelShort, dataPortalName,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM series LEFT JOIN geographies ON name LIKE CONCAT('%@', handle, '.%')
	WHERE
	(SELECT list FROM data_lists JOIN categories WHERE categories.data_list_id = data_lists.id AND categories.id = ?)
	REGEXP CONCAT('[[:<:]]', TRIM(TRAILING 'NS' FROM left(name, locate('@', name) - 1)), '(NS){0,1}@.*[[:>:]]')`
var geoFilter = ` AND series.name LIKE CONCAT('%@', ? ,'.%') `
var freqFilter = ` AND series.name LIKE CONCAT('%@%.', ?) `
var sortStmt = ` ORDER BY LOCATE(CONCAT(TRIM(TRAILING 'NS' FROM left(name, locate('@', name) - 1)), '@'),
	(SELECT list FROM data_lists JOIN categories WHERE categories.data_list_id = data_lists.id AND categories.id = ?)) +
	LOCATE(CONCAT(TRIM(TRAILING 'NS' FROM left(name, locate('@', name) - 1)), 'NS@'),
	(SELECT list FROM data_lists JOIN categories WHERE categories.data_list_id = data_lists.id AND categories.id = ?));`
var siblingsPrefix = `SELECT series.id, series.name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM series LEFT JOIN geographies ON name LIKE CONCAT('%@', handle, '.%')
	JOIN (SELECT name FROM series where id = ?) as original_series
	WHERE series.name LIKE CONCAT(left(original_series.name, locate("@", original_series.name)), '%')
	AND series.name LIKE CONCAT('%@', ?, '.', ?);`

func (r *SeriesRepository) GetSeriesByCategoryAndFreq(
	categoryId int64,
	freq string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(
		strings.Join([]string{seriesPrefix, freqFilter, sortStmt}, ""),
		categoryId,
		freq,
		categoryId,
		categoryId,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr  := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesByCategoryGeoAndFreq(
	categoryId int64,
	geoHandle string,
	freq string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(
		strings.Join([]string{seriesPrefix, geoFilter, freqFilter, sortStmt}, ""),
		categoryId,
		geoHandle,
		freq,
		categoryId,
		categoryId,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr  := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesByCategoryAndGeo(
	categoryId int64,
	geoHandle string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(
		strings.Join([]string{seriesPrefix, geoFilter, sortStmt}, ""),
		categoryId,
		geoHandle,
		categoryId,
		categoryId,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr  := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesBySearchText(searchText string) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT series.id, name, description, frequency,
	!(name REGEXP '.*NS@.*') AS seasonally_adjusted,
	unitsLabel, unitsLabelShort, dataPortalName,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM series LEFT JOIN geographies ON name LIKE CONCAT('%@', handle, '.%')
	WHERE
	((MATCH(name, description, dataPortalName)
	AGAINST(? IN NATURAL LANGUAGE MODE))
	OR LOWER(CONCAT(name, description, dataPortalName)) LIKE CONCAT('%', LOWER(?), '%'));`, searchText, searchText)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr  := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesByCategory(categoryId int64) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(
		strings.Join([]string{seriesPrefix, sortStmt}, ""),
		categoryId,
		categoryId,
		categoryId,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		dataPortalSeries, scanErr  := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesSiblingsById(seriesId int64) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT series.id, series.name, description, frequency,
	!(name REGEXP '.*NS@.*') AS seasonally_adjusted,
	unitsLabel, unitsLabelShort, dataPortalName,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM series LEFT JOIN geographies ON name LIKE CONCAT('%@', handle, '.%')
	JOIN (SELECT name FROM series where id = ?) as original_series
	WHERE series.name LIKE CONCAT(left(original_series.name, locate("@", original_series.name)), '%');`, seriesId)
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

func (r *SeriesRepository) GetSeriesSiblingsByIdAndFreq(seriesId int64, freq string) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT series.id, series.name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM series LEFT JOIN geographies ON name LIKE CONCAT('%@', handle, '.%')
	JOIN (SELECT name FROM series where id = ?) as original_series
	WHERE series.name LIKE CONCAT(left(original_series.name, locate("@", original_series.name)), '%')
	AND series.name LIKE CONCAT('%@%.', ?);`, seriesId, freq)
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

func (r *SeriesRepository) GetSeriesSiblingsByIdAndGeo(seriesId int64, geo string) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT series.id, series.name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM series LEFT JOIN geographies ON name LIKE CONCAT('%@', handle, '.%')
	JOIN (SELECT name FROM series where id = ?) as original_series
	WHERE series.name LIKE CONCAT(left(original_series.name, locate("@", original_series.name)), '%')
	AND series.name LIKE CONCAT('%@%', ? ,'%.%');`, seriesId, geo)
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

func (r *SeriesRepository) GetSeriesSiblingsByIdGeoAndFreq(seriesId int64, geo string, freq string) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT series.id, series.name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM series LEFT JOIN geographies ON name LIKE CONCAT('%@', handle, '.%')
	JOIN (SELECT name FROM series where id = ?) as original_series
	WHERE series.name LIKE CONCAT(left(original_series.name, locate("@", original_series.name)), '%')
	AND series.name LIKE CONCAT('%@', ?, '.', ?);`, seriesId, geo, freq)
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

func (r *SeriesRepository) GetSeriesSiblingsFreqById(seriesId int64) (frequencyList []models.FrequencyResult, err error) {
	rows, err := r.DB.Query(`SELECT DISTINCT(RIGHT(series.name, 1)) as freq, frequency
	FROM series JOIN (SELECT name FROM series where id = ?) as original_series
	WHERE series.name LIKE CONCAT(left(original_series.name, locate("@", original_series.name)), '%');`, seriesId)
	if err != nil {
		return
	}
	for rows.Next() {
		frequency := models.Frequency{}
		err = rows.Scan(
			&frequency.Freq,
			&frequency.Label,
		)
		if err != nil {
			return
		}
		frequencyResult := models.FrequencyResult{Freq: frequency.Freq}
		if frequency.Label.Valid {
			frequencyResult.Label = frequency.Label.String
		}
		frequencyList = append(frequencyList, frequencyResult)
	}
	return
}

func (r *SeriesRepository) GetSeriesById(seriesId int64) (dataPortalSeries models.DataPortalSeries, err error) {
	row := r.DB.QueryRow(`SELECT series.id, name, description, frequency,
	!(name REGEXP '.*NS@.*') AS seasonally_adjusted,
	unitsLabel, unitsLabelShort, dataPortalName,
	fips, SUBSTRING_INDEX(SUBSTR(series.name, LOCATE('@', series.name) + 1), '.', 1) as shandle, display_name_short
	FROM series LEFT JOIN geographies ON name LIKE CONCAT('%@', handle, '.%')
	WHERE series.id = ?;`, seriesId)
	dataPortalSeries, err = getNextSeriesFromRow(row)
	return
}

func (r *SeriesRepository) GetSeriesObservations(seriesId int64) (seriesObservations models.SeriesObservations, err error) {
	lvlTransform, start, end, err := r.GetTransformation(Levels, seriesId)
	if err != nil {
		return
	}
	yoyTransform, yoyStart, yoyEnd, err := r.GetTransformation(YoyPCh, seriesId)
	if err != nil {
		return
	}
	if yoyStart.Before(start) {
		start = yoyStart
	}
	if end.Before(yoyEnd) {
		end = yoyEnd
	}
	ytdTransform, ytdStart, ytdEnd, err := r.GetTransformation(YTD, seriesId)
	if err != nil {
		return
	}
	if ytdStart.Before(start) {
		start = ytdStart
	}
	if end.Before(ytdEnd) {
		end = ytdEnd
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

func (r *SeriesRepository) GetTransformation(transformation string, seriesId int64) (
	transformationResult models.TransformationResult,
	observationStart time.Time,
	observationEnd time.Time,
	err error,
) {
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
		observations = append(
			observations,
			models.DataPortalObservation{
				Date:  observation.Date,
				Value: observation.Value.Float64,
			},
		)
	}
	transformationResult.Transformation = transformation
	transformationResult.Observations = observations
	return
}
