package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"time"
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
}

func (r *SeriesRepository) GetSeriesByCategoryAndFreq(
	categoryId int64,
	freq string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(
		`SELECT id, name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName
	FROM series WHERE
	(SELECT list FROM data_lists JOIN categories WHERE categories.data_list_id = data_lists.id AND categories.id = ?)
	REGEXP CONCAT('[[:<:]]', left(name, locate("@", name)), '.*[[:>:]]') AND
	name LIKE CONCAT('%@%.', ?);`,
		categoryId,
		freq,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		series := models.Series{}
		err = rows.Scan(
			&series.Id,
			&series.Name,
			&series.Description,
			&series.Frequency,
			&series.SeasonallyAdjusted,
			&series.UnitsLabel,
			&series.UnitsLabelShort,
			&series.DataPortalName,
		)
		if err != nil {
			return
		}
		dataPortalSeries := models.DataPortalSeries{Id: series.Id, Name: series.Name}
		if series.DataPortalName.Valid {
			dataPortalSeries.Title = series.DataPortalName.String
		}
		if series.Description.Valid {
			dataPortalSeries.Description = series.Description.String
		}
		if series.Frequency.Valid {
			dataPortalSeries.Frequency = series.Frequency.String
		}
		if series.SeasonallyAdjusted.Valid {
			dataPortalSeries.SeasonallyAdjusted = series.SeasonallyAdjusted.Bool
		}
		if series.UnitsLabel.Valid {
			dataPortalSeries.UnitsLabel = series.UnitsLabel.String
		}
		if series.UnitsLabelShort.Valid {
			dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
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
		`SELECT id, name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName
	FROM series WHERE
	(SELECT list FROM data_lists JOIN categories WHERE categories.data_list_id = data_lists.id AND categories.id = ?)
	REGEXP CONCAT('[[:<:]]', left(name, locate("@", name)), '.*[[:>:]]') AND name LIKE CONCAT('%@%', ? ,'%.%') AND
	name LIKE CONCAT('%@%.', ?);`,
		categoryId,
		geoHandle,
		freq,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		series := models.Series{}
		err = rows.Scan(
			&series.Id,
			&series.Name,
			&series.Description,
			&series.Frequency,
			&series.SeasonallyAdjusted,
			&series.UnitsLabel,
			&series.UnitsLabelShort,
			&series.DataPortalName,
		)
		if err != nil {
			return
		}
		dataPortalSeries := models.DataPortalSeries{Id: series.Id, Name: series.Name}
		if series.DataPortalName.Valid {
			dataPortalSeries.Title = series.DataPortalName.String
		}
		if series.Description.Valid {
			dataPortalSeries.Description = series.Description.String
		}
		if series.Frequency.Valid {
			dataPortalSeries.Frequency = series.Frequency.String
		}
		if series.SeasonallyAdjusted.Valid {
			dataPortalSeries.SeasonallyAdjusted = series.SeasonallyAdjusted.Bool
		}
		if series.UnitsLabel.Valid {
			dataPortalSeries.UnitsLabel = series.UnitsLabel.String
		}
		if series.UnitsLabelShort.Valid {
			dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
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
		`SELECT id, name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName
	FROM series WHERE
	(SELECT list FROM data_lists JOIN categories WHERE categories.data_list_id = data_lists.id AND categories.id = ?)
	REGEXP CONCAT('[[:<:]]', left(name, locate("@", name)), '.*[[:>:]]') AND name LIKE CONCAT('%@%', ? ,'%.%');`,
		categoryId,
		geoHandle,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		series := models.Series{}
		err = rows.Scan(
			&series.Id,
			&series.Name,
			&series.Description,
			&series.Frequency,
			&series.SeasonallyAdjusted,
			&series.UnitsLabel,
			&series.UnitsLabelShort,
			&series.DataPortalName,
		)
		if err != nil {
			return
		}
		dataPortalSeries := models.DataPortalSeries{Id: series.Id, Name: series.Name}
		if series.DataPortalName.Valid {
			dataPortalSeries.Title = series.DataPortalName.String
		}
		if series.Description.Valid {
			dataPortalSeries.Description = series.Description.String
		}
		if series.Frequency.Valid {
			dataPortalSeries.Frequency = series.Frequency.String
		}
		if series.SeasonallyAdjusted.Valid {
			dataPortalSeries.SeasonallyAdjusted = series.SeasonallyAdjusted.Bool
		}
		if series.UnitsLabel.Valid {
			dataPortalSeries.UnitsLabel = series.UnitsLabel.String
		}
		if series.UnitsLabelShort.Valid {
			dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesByCategory(categoryId int64) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT id, name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName
	FROM series WHERE
	(SELECT list FROM data_lists JOIN categories WHERE categories.data_list_id = data_lists.id AND categories.id = ?)
	REGEXP CONCAT('[[:<:]]', left(name, locate("@", name)), '.*[[:>:]]');`, categoryId)
	if err != nil {
		return
	}
	for rows.Next() {
		series := models.Series{}
		err = rows.Scan(
			&series.Id,
			&series.Name,
			&series.Description,
			&series.Frequency,
			&series.SeasonallyAdjusted,
			&series.UnitsLabel,
			&series.UnitsLabelShort,
			&series.DataPortalName,
		)
		if err != nil {
			return
		}
		dataPortalSeries := models.DataPortalSeries{Id: series.Id, Name: series.Name}
		if series.DataPortalName.Valid {
			dataPortalSeries.Title = series.DataPortalName.String
		}
		if series.Description.Valid {
			dataPortalSeries.Description = series.Description.String
		}
		if series.Frequency.Valid {
			dataPortalSeries.Frequency = series.Frequency.String
		}
		if series.SeasonallyAdjusted.Valid {
			dataPortalSeries.SeasonallyAdjusted = series.SeasonallyAdjusted.Bool
		}
		if series.UnitsLabel.Valid {
			dataPortalSeries.UnitsLabel = series.UnitsLabel.String
		}
		if series.UnitsLabelShort.Valid {
			dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesSiblingsById(seriesId int64) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT id, series.name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName
	FROM series JOIN (SELECT name FROM series where id = ?) as original_series
	WHERE series.name LIKE CONCAT(left(original_series.name, locate("@", original_series.name)), '%');`, seriesId)
	if err != nil {
		return
	}
	for rows.Next() {
		series := models.Series{}
		err = rows.Scan(
			&series.Id,
			&series.Name,
			&series.Description,
			&series.Frequency,
			&series.SeasonallyAdjusted,
			&series.UnitsLabel,
			&series.UnitsLabelShort,
			&series.DataPortalName,
		)
		if err != nil {
			return
		}
		dataPortalSeries := models.DataPortalSeries{Id: series.Id, Name: series.Name}
		if series.DataPortalName.Valid {
			dataPortalSeries.Title = series.DataPortalName.String
		}
		if series.Description.Valid {
			dataPortalSeries.Description = series.Description.String
		}
		if series.Frequency.Valid {
			dataPortalSeries.Frequency = series.Frequency.String
		}
		if series.SeasonallyAdjusted.Valid {
			dataPortalSeries.SeasonallyAdjusted = series.SeasonallyAdjusted.Bool
		}
		if series.UnitsLabel.Valid {
			dataPortalSeries.UnitsLabel = series.UnitsLabel.String
		}
		if series.UnitsLabelShort.Valid {
			dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesSiblingsByIdAndFreq(seriesId int64, freq string) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT id, series.name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName
	FROM series JOIN (SELECT name FROM series where id = ?) as original_series
	WHERE series.name LIKE CONCAT(left(original_series.name, locate("@", original_series.name)), '%')
	AND series.name LIKE CONCAT('%@%.', ?);`, seriesId, freq)
	if err != nil {
		return
	}
	for rows.Next() {
		series := models.Series{}
		err = rows.Scan(
			&series.Id,
			&series.Name,
			&series.Description,
			&series.Frequency,
			&series.SeasonallyAdjusted,
			&series.UnitsLabel,
			&series.UnitsLabelShort,
			&series.DataPortalName,
		)
		if err != nil {
			return
		}
		dataPortalSeries := models.DataPortalSeries{Id: series.Id, Name: series.Name}
		if series.DataPortalName.Valid {
			dataPortalSeries.Title = series.DataPortalName.String
		}
		if series.Description.Valid {
			dataPortalSeries.Description = series.Description.String
		}
		if series.Frequency.Valid {
			dataPortalSeries.Frequency = series.Frequency.String
		}
		if series.SeasonallyAdjusted.Valid {
			dataPortalSeries.SeasonallyAdjusted = series.SeasonallyAdjusted.Bool
		}
		if series.UnitsLabel.Valid {
			dataPortalSeries.UnitsLabel = series.UnitsLabel.String
		}
		if series.UnitsLabelShort.Valid {
			dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesSiblingsByIdAndGeo(seriesId int64, geo string) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT id, series.name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName
	FROM series JOIN (SELECT name FROM series where id = ?) as original_series
	WHERE series.name LIKE CONCAT(left(original_series.name, locate("@", original_series.name)), '%')
	AND series.name LIKE CONCAT('%@%', ? ,'%.%');`, seriesId, geo)
	if err != nil {
		return
	}
	for rows.Next() {
		series := models.Series{}
		err = rows.Scan(
			&series.Id,
			&series.Name,
			&series.Description,
			&series.Frequency,
			&series.SeasonallyAdjusted,
			&series.UnitsLabel,
			&series.UnitsLabelShort,
			&series.DataPortalName,
		)
		if err != nil {
			return
		}
		dataPortalSeries := models.DataPortalSeries{Id: series.Id, Name: series.Name}
		if series.DataPortalName.Valid {
			dataPortalSeries.Title = series.DataPortalName.String
		}
		if series.Description.Valid {
			dataPortalSeries.Description = series.Description.String
		}
		if series.Frequency.Valid {
			dataPortalSeries.Frequency = series.Frequency.String
		}
		if series.SeasonallyAdjusted.Valid {
			dataPortalSeries.SeasonallyAdjusted = series.SeasonallyAdjusted.Bool
		}
		if series.UnitsLabel.Valid {
			dataPortalSeries.UnitsLabel = series.UnitsLabel.String
		}
		if series.UnitsLabelShort.Valid {
			dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *SeriesRepository) GetSeriesSiblingsByIdGeoAndFreq(seriesId int64, geo string, freq string) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT id, series.name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName
	FROM series JOIN (SELECT name FROM series where id = ?) as original_series
	WHERE series.name LIKE CONCAT(left(original_series.name, locate("@", original_series.name)), '%')
	AND series.name LIKE CONCAT('%@', ?, '.', ?);`, seriesId, geo, freq)
	if err != nil {
		return
	}
	for rows.Next() {
		series := models.Series{}
		err = rows.Scan(
			&series.Id,
			&series.Name,
			&series.Description,
			&series.Frequency,
			&series.SeasonallyAdjusted,
			&series.UnitsLabel,
			&series.UnitsLabelShort,
			&series.DataPortalName,
		)
		if err != nil {
			return
		}
		dataPortalSeries := models.DataPortalSeries{Id: series.Id, Name: series.Name}
		if series.DataPortalName.Valid {
			dataPortalSeries.Title = series.DataPortalName.String
		}
		if series.Description.Valid {
			dataPortalSeries.Description = series.Description.String
		}
		if series.Frequency.Valid {
			dataPortalSeries.Frequency = series.Frequency.String
		}
		if series.SeasonallyAdjusted.Valid {
			dataPortalSeries.SeasonallyAdjusted = series.SeasonallyAdjusted.Bool
		}
		if series.UnitsLabel.Valid {
			dataPortalSeries.UnitsLabel = series.UnitsLabel.String
		}
		if series.UnitsLabelShort.Valid {
			dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
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
	row := r.DB.QueryRow(`SELECT id, name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName
	FROM series WHERE id = ?`, seriesId)
	series := models.Series{}
	err = row.Scan(
		&series.Id,
		&series.Name,
		&series.Description,
		&series.Frequency,
		&series.SeasonallyAdjusted,
		&series.UnitsLabel,
		&series.UnitsLabelShort,
		&series.DataPortalName,
	)
	if err != nil {
		return
	}
	dataPortalSeries = models.DataPortalSeries{Id: series.Id, Name: series.Name}
	dataPortalSeries.FrequencyShort = series.Name[len(series.Name)-1:]
	if series.DataPortalName.Valid {
		dataPortalSeries.Title = series.DataPortalName.String
	}
	if series.Description.Valid {
		dataPortalSeries.Description = series.Description.String
	}
	if series.Frequency.Valid {
		dataPortalSeries.Frequency = series.Frequency.String
	}
	if series.SeasonallyAdjusted.Valid {
		dataPortalSeries.SeasonallyAdjusted = series.SeasonallyAdjusted.Bool
	}
	if series.UnitsLabel.Valid {
		dataPortalSeries.UnitsLabel = series.UnitsLabel.String
	}
	if series.UnitsLabelShort.Valid {
		dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
	}
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
	seriesObservations.TransformationResults = []models.TransformationResult{lvlTransform, yoyTransform}
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
