package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"time"
)

type SeriesRepository struct {
	DB *sql.DB
}

const (
	Levels = "lvl"
	YoyPCh = "pc1"
)

var transformations map[string]string = map[string]string{
	"lvl": `SELECT date, value FROM data_points WHERE series_id = ? and current = 1;`,
	"pc1": `SELECT t1.date, (t1.value/t2.last_value - 1)*100 AS yoy
		FROM (SELECT value, date, DATE_SUB(date, INTERVAL 1 YEAR) AS last_year
		FROM data_points WHERE series_id = ? AND current = 1) AS t1
		LEFT JOIN (SELECT value AS last_value, date
		FROM data_points WHERE series_id = 146634 and current = 1) AS t2
		ON (t1.last_year = t2.date);`,
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
	LIKE CONCAT('%', left(name, locate("@", name)), '%') AND name LIKE CONCAT('%@%', ? ,'%.%');`,
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
	LIKE CONCAT('%', left(name, locate("@", name)), '%');`, categoryId)
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

func (r *SeriesRepository) GetTransformation(transformation string, seriesId int64) (
	transformationResult models.TransformationResult,
	observationStart time.Time,
	observationEnd time.Time,
	err error,
) {
	rows, err := r.DB.Query(transformations[transformation], seriesId)
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
				Date: observation.Date,
				Value: observation.Value.Float64,
			},
		)
	}
	transformationResult.Transformation = transformation
	transformationResult.Observations = observations
	return
}
