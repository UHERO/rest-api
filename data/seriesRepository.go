package data

import (
	"database/sql"
	"strings"
	"time"

	"github.com/UHERO/rest-api/models"
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
	MOMPercentChange  = "mompc1"
	YOYChange         = "ch1"
	YTDChange         = "ytdch1"
	C5MAChange        = "c5mach1"
	MOMChange         = "momch1"
)

var transformations = map[string]transformation{
	Levels: { // untransformed value
		//language=MySQL
		Statement: `SELECT date, value, (pseudo_history = b'1'), series.decimals
					FROM <%DATAPOINTS%> dp
					LEFT JOIN <%SERIES%> AS series ON series.id = dp.series_id
					WHERE dp.series_id = ? `,
		PlaceholderCount: 1,
		Label:            "lvl",
	},
	YOYPercentChange: { // percent change from 1 year ago
		//language=MySQL
		Statement: `SELECT t1.date, (t1.value/t2.value - 1) * 100 AS yoy,
							(t1.pseudo_history = true AND t2.pseudo_history = true) AS ph, series.decimals
					FROM <%DATAPOINTS%> AS t1
					LEFT JOIN <%DATAPOINTS%> AS t2 ON t2.series_id = t1.series_id
										  		  AND t2.date = DATE_SUB(t1.date, INTERVAL 1 YEAR)
					JOIN <%SERIES%> AS series ON series.id = t1.series_id
					WHERE t1.series_id = ? `,
		PlaceholderCount: 1,
		Label:            "pc1",
	},

	YOYChange: { // change from 1 year ago
		//language=MySQL
		Statement: `SELECT t1.date, (t1.value - t2.value) AS yoy,
			 			(t1.pseudo_history = true AND t2.pseudo_history = true) AS ph, series.decimals
					FROM <%DATAPOINTS%> AS t1
					LEFT JOIN <%DATAPOINTS%> AS t2 ON t2.series_id = t1.series_id
										 		  AND t2.date = DATE_SUB(t1.date, INTERVAL 1 YEAR)
					JOIN <%SERIES%> AS series ON series.id = t1.series_id
					WHERE t1.series_id = ? `,
		PlaceholderCount: 1,
		Label:            "pc1",
	},

	YTDChange: { // ytd change from 1 year ago
		//language=MySQL
		Statement: `
	    WITH t1 AS (
			SELECT date,
				@sum := IF(@yr = year(date), @sum, 0) + value AS ytd_sum,
				@count := IF(@yr = year(date), @count, 0) + 1 AS count,
				@yr := year(date),
			    series_id,
			    pseudo_history
			FROM <%DATAPOINTS%>
			JOIN (SELECT @sum := null, @yr := null) AS init
			WHERE series_id = ?
		), t2 AS (
			SELECT date, @sum := IF(@yr = year(date), @sum, 0) + value AS ytd_sum,
				@count := IF(@yr = year(date), @count, 0) + 1 AS count,
				@yr := year(date),
			    pseudo_history
			FROM <%DATAPOINTS%>
			JOIN (SELECT @sum := null, @yr := null) AS init
			WHERE series_id = ?
		)
		SELECT t1.date, (t1.ytd_sum / t1.count - t2.ytd_sum / t2.count) AS ytd_change,
				(t1.pseudo_history = true) AND (t2.pseudo_history = true) AS ph,
				series.decimals
		FROM t1
		LEFT JOIN t2 ON t2.date = date_sub(t1.date, INTERVAL 1 YEAR)
		JOIN <%SERIES%> AS series ON t1.series_id = series.id `,
		PlaceholderCount: 2,
		Label:            "ytd",
	},

	YTDPercentChange: { // ytd percent change from 1 year ago
		//language=MySQL
		Statement: `
		WITH t1 AS (
			SELECT date,
				   @sum := if(@yr = year(date), @sum, 0) + value AS ytd_sum,
				   @yr := year(date),
				   series_id,
				   pseudo_history
			FROM <%DATAPOINTS%>
			JOIN (SELECT @sum := null, @yr := null) AS init
			WHERE series_id = ?
		), t2 AS (
			SELECT date,
				   @sum := if(@yr = year(date), @sum, 0) + value AS ytd_sum,
				   @yr := year(date),
				   pseudo_history
			FROM <%DATAPOINTS%>
			JOIN (SELECT @sum := null, @yr := null) AS init
			WHERE series_id = ?
		)
		SELECT t1.date, (t1.ytd_sum / t2.ytd_sum - 1) * 100 AS ytd_pct_change,
			(t1.pseudo_history = true AND t2.pseudo_history = true) AS ph,
		    series.decimals
		FROM t1
		LEFT JOIN t2 ON t2.date = date_sub(t1.date, INTERVAL 1 YEAR)
		JOIN <%SERIES%> AS series ON series.id = t1.series_id `,
		PlaceholderCount: 2,
		Label:            "ytd",
	},

	C5MAPercentChange: { // c5ma percent change from 1 year ago
		//language=MySQL
		Statement: `
		WITH c5ma_agg AS (
			SELECT p1.series_id, p1.date, p1.pseudo_history, CASE WHEN count(*) = 5 THEN AVG(p2.value) ELSE NULL END AS c5ma
			FROM <%DATAPOINTS%> AS p1
			JOIN <%DATAPOINTS%> AS p2 ON p2.series_id = p1.series_id
								     AND p2.date BETWEEN DATE_SUB(p1.date, INTERVAL 2 YEAR)
													 AND DATE_ADD(p1.date, INTERVAL 2 YEAR)
			WHERE p1.series_id = ?
			GROUP BY 1, 2, 3
		)
		SELECT cur.date, (cur.c5ma / lastyear.c5ma - 1) * 100 AS c5ma_pct_change,
			  (cur.pseudo_history = true AND lastyear.pseudo_history = true) AS ph, series.decimals
		FROM c5ma_agg AS cur
		JOIN c5ma_agg AS lastyear ON lastyear.date = DATE_SUB(cur.date, INTERVAL 1 YEAR)
		JOIN <%SERIES%> AS series ON series.id = cur.series_id `,
		PlaceholderCount: 1,
		Label:            "c5ma",
	},
	C5MAChange: { // cm5a change from 1 year ago
		//language=MySQL
		Statement: `
		WITH c5ma_agg AS (
			SELECT p1.series_id, p1.date, p1.pseudo_history, CASE WHEN count(*) = 5 THEN AVG(p2.value) ELSE NULL END AS c5ma
			FROM <%DATAPOINTS%> AS p1
			JOIN <%DATAPOINTS%> AS p2 ON p2.series_id = p1.series_id
								     AND p2.date BETWEEN DATE_SUB(p1.date, INTERVAL 2 YEAR)
													 AND DATE_ADD(p1.date, INTERVAL 2 YEAR)
			WHERE p1.series_id = ?
			GROUP BY 1, 2, 3
		)
		SELECT cur.date, (cur.c5ma - lastyear.c5ma) AS c5ma_change,
			  (cur.pseudo_history = true AND lastyear.pseudo_history = true) AS ph, series.decimals
		FROM c5ma_agg AS cur
		JOIN c5ma_agg AS lastyear ON lastyear.date = DATE_SUB(cur.date, INTERVAL 1 YEAR)
		JOIN <%SERIES%> AS series ON series.id = cur.series_id `,
		PlaceholderCount: 1,
		Label:            "c5ma",
	},

	MOMPercentChange: { // mom percent change from 1 month ago
		//language=MySQL
		Statement: `
			WITH t1 AS (
			    SELECT date, sum(value) AS month_sum, series_id, pseudo_history
				FROM <%DATAPOINTS%>
			 	WHERE series_id = ?
				GROUP BY YEAR(date), MONTH(date)
			)
			SELECT date, (month_sum / LAG(month_sum, 1) OVER(ORDER BY date) - 1) * 100 AS mom_pct_change,
			       t1.pseudo_history = true AS ph, series.decimals
			FROM t1
			JOIN <%SERIES%> AS series ON series.id = t1.series_id`,
		PlaceholderCount: 1,
		Label:            "mom",
	},

	MOMChange: { // mom change from 1 month ago
		//language=MySQL
		Statement: `
			WITH t1 AS (
			    SELECT date, sum(value) AS month_sum, series_id, pseudo_history
				FROM <%DATAPOINTS%>
			    WHERE series_id = ?
				GROUP BY YEAR(date), MONTH(date)
			)
			SELECT date, (month_sum - LAG(month_sum, 1) OVER(ORDER BY date) - 1) AS mom_change,
			       t1.pseudo_history = true AS ph, series.decimals
			FROM t1
			JOIN <%SERIES%> AS series ON series.id = t1.series_id`,
		PlaceholderCount: 1,
		Label:            "mom",
	},
}

//language=MySQL
var seriesPrefix = `
	SELECT series_id, series_name, series_universe, series_description, frequency, seasonally_adjusted, seasonal_adjustment,
	       units_long, units_short, data_portal_name, percent, pv.real, source_description, source_link, source_detail_description,
		   MAX(table_prefix), MAX(table_postfix), MAX(measurement_id), MAX(measurement_portal_name), MAX(dlm_indent),
	       base_year, decimals, geo_fips, geo_handle, geo_display_name, geo_display_name_short
	FROM <%PORTAL%> pv
	WHERE category_id = ? `

//language=MySQL
var measurementSeriesPrefix = `
	SELECT DISTINCT series_id, series_name, series_universe, series_description, frequency, seasonally_adjusted, seasonal_adjustment,
	       units_long, units_short, data_portal_name, percent, pv.real, source_description, source_link, source_detail_description,
		   table_prefix, table_postfix, measurement_id, measurement_portal_name, NULL,
	       base_year, decimals, geo_fips, geo_handle, geo_display_name, geo_display_name_short
	FROM <%PORTAL%> pv
	WHERE measurement_id = ? `
var fcFilter = ` AND series_name REGEXP ? `
var geoFilter = ` AND geo_handle = ? `
var freqFilter = ` AND frequency = ? `
var measurementPostfix = " ;" // this part of query no longer needed, but too troublesome to change all the code
var sortStmt = ` GROUP BY series_id ORDER BY MAX(dlm_list_order);`
var siblingSortStmt = ` GROUP BY series_id;`

//language=MySQL
var siblingsPrefix = `
	SELECT series_id, series_name, series_universe, series_description, frequency, seasonally_adjusted, seasonal_adjustment,
	       units_long, units_short, data_portal_name, percent, pv.real, source_description, source_link, source_detail_description,
		   MAX(table_prefix), MAX(table_postfix), MAX(pv.measurement_id), MAX(measurement_portal_name), MAX(dlm_indent),
	       base_year, decimals, geo_fips, geo_handle, geo_display_name, geo_display_name_short
	FROM <%PORTAL%> pv
	JOIN (SELECT measurement_id FROM measurement_series WHERE series_id = ?) AS mfilt ON mfilt.measurement_id = pv.measurement_id
	WHERE EXISTS(SELECT * FROM public_data_points WHERE series_id = pv.series_id) `

func (r *FooRepository) GetSeriesByGroupAndFreq(
	groupId int64,
	freq string,
	groupType GroupType,
) (seriesList []models.DataPortalSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	catId := groupId
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = measurementPostfix
		catId = 0
	}
	rows, err := r.RunQuery(
		strings.Join([]string{prefix, freqFilter, sort}, ""),
		groupId,
		freqDbNames[strings.ToUpper(freq)],
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id, catId)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *FooRepository) GetSeriesByGroupGeoAndFreq(
	groupId int64,
	geoHandle string,
	freq string,
	forecast string,
	groupType GroupType,
) (seriesList []models.DataPortalSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	catId := groupId
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = measurementPostfix
		catId = 0
	}
	rows, err := r.RunQuery(
		strings.Join([]string{prefix, fcFilter, geoFilter, freqFilter, sort}, ""),
		groupId,
		forecast,
		geoHandle,
		freqDbNames[strings.ToUpper(freq)],
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id, catId)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *FooRepository) GetInflatedSeriesByGroupGeoAndFreq(
	groupId int64,
	geoHandle string,
	freq string,
	forecast string,
	groupType GroupType,
) (seriesList []models.InflatedSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	catId := groupId
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = measurementPostfix
		catId = 0
	}
	rows, err := r.RunQuery(
		strings.Join([]string{prefix, fcFilter, geoFilter, freqFilter, sort}, ""),
		groupId,
		forecast,
		geoHandle,
		freqDbNames[strings.ToUpper(freq)],
	)
	if err != nil {
		return
	}
	defer rows.Close()

	seriesList = make([]models.InflatedSeries, 0, 30)
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id, catId)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesObservations, scanErr := r.GetSeriesObservations(dataPortalSeries.Id, "")
		if scanErr != nil {
			return seriesList, scanErr
		}
		seriesList = append(seriesList, models.InflatedSeries{dataPortalSeries, seriesObservations})
	}
	return
}

func (r *FooRepository) GetSeriesByGroupAndGeo(
	groupId int64,
	geoHandle string,
	groupType GroupType,
) (seriesList []models.DataPortalSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	catId := groupId
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = measurementPostfix
		catId = 0
	}
	rows, err := r.RunQuery(
		strings.Join([]string{prefix, geoFilter, sort}, ""),
		groupId,
		geoHandle,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id, catId)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *FooRepository) GetInflatedSeriesByGroup(
	groupId int64,
	groupType GroupType,
) (seriesList []models.InflatedSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	catId := groupId
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = measurementPostfix
		catId = 0
	}
	rows, err := r.RunQuery(
		strings.Join([]string{prefix, sort}, ""),
		groupId,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id, catId)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesObservations, scanErr := r.GetSeriesObservations(dataPortalSeries.Id, "")
		if scanErr != nil {
			return seriesList, scanErr
		}
		inflatedSeries := models.InflatedSeries{dataPortalSeries, seriesObservations}
		seriesList = append(seriesList, inflatedSeries)
	}
	return
}

func (r *FooRepository) GetSeriesByGroup(
	groupId int64,
	groupType GroupType,
) (seriesList []models.DataPortalSeries, err error) {
	prefix := seriesPrefix
	sort := sortStmt
	catId := groupId
	if groupType == Measurement {
		prefix = measurementSeriesPrefix
		sort = measurementPostfix
		catId = 0
	}
	rows, err := r.RunQuery(
		strings.Join([]string{prefix, sort}, ""),
		groupId,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id, catId)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *FooRepository) GetFreqByCategory(categoryId int64, forecast string) (frequencies []models.DataPortalFrequency, err error) {
	//language=MySQL
	rows, err := r.RunQuery(`SELECT DISTINCT(RIGHT(series_name, 1)) AS freq
							 FROM <%PORTAL%> pv
							 WHERE category_id = ?
							 AND series_name REGEXP ?
							 ORDER BY FIELD(freq, "A", "S", "Q", "M", "W", "D");`, categoryId, forecast)
	if err != nil {
		return
	}
	defer rows.Close()
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

func (r *FooRepository) GetForecastByCategory(categoryId int64) (forecasts []string, err error) {
	//language=MySQL
	rows, err := r.RunQuery(`SELECT DISTINCT SUBSTRING_INDEX(SUBSTRING_INDEX(series_name, '@', 1), '&', -1)
							 FROM <%PORTAL%> pv
							 WHERE category_id = ? ORDER BY 1;`, categoryId)
	if err != nil {
		return
	}
	defer rows.Close()
	var fcName string
	for rows.Next() {
		err = rows.Scan(&fcName)
		if err != nil {
			return
		}
		forecasts = append(forecasts, fcName)
	}
	return
}

func (r *FooRepository) GetForecastsById(seriesId int64) (forecasts []models.DataPortalForecast, err error) {
	//language=MySQL
	rows, err := r.RunQuery(`SELECT DISTINCT SUBSTRING_INDEX(SUBSTRING_INDEX(series_name, '@', 1), '&', -1) AS fc, pv.frequency
							 FROM <%PORTAL%> pv
							 JOIN (SELECT * FROM <%SERIES%> WHERE id = ?) s
								 ON SUBSTRING_INDEX(s.name, '&', 1) = SUBSTRING_INDEX(pv.series_name, '&', 1)
								AND s.geography_id = pv.series_geo_id
							 WHERE pv.series_universe = 'FC' ORDER BY 1;`, seriesId)
	if err != nil {
		return
	}
	defer rows.Close()
	var fcName, freqFull string
	for rows.Next() {
		err = rows.Scan(&fcName, &freqFull)
		if err != nil {
			return
		}
		code := freqDbCodes[freqFull]
		forecasts = append(forecasts, models.DataPortalForecast{Forecast: fcName, Freq: code, Label: freqLabel[code]})
	}
	return
}

func (r *FooRepository) GetSeriesSiblingsById(seriesId int64, forecast string, categoryId int64) (seriesList []models.DataPortalSeries, err error) {
	foofar := fcFilter
	if forecast != "@" {
		foofar += " AND "
	}
	rows, err := r.RunQuery(
		strings.Join([]string{siblingsPrefix, foofar, siblingSortStmt}, ""),
		seriesId,
		forecast,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id, categoryId)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *FooRepository) GetSeriesSiblingsByIdAndFreq(
	seriesId int64,
	freq string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.RunQuery(strings.Join(
		[]string{siblingsPrefix, freqFilter, siblingSortStmt}, ""),
		seriesId,
		freqDbNames[strings.ToUpper(freq)],
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id, 0)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *FooRepository) GetSeriesSiblingsByIdAndGeo(
	seriesId int64,
	geo string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.RunQuery(
		strings.Join([]string{siblingsPrefix, geoFilter, siblingSortStmt}, ""),
		seriesId,
		geo,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id, 0)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *FooRepository) GetSeriesSiblingsByIdGeoAndFreq(
	seriesId int64,
	geo string,
	freq string,
) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.RunQuery(
		strings.Join([]string{siblingsPrefix, geoFilter, freqFilter, siblingSortStmt}, ""),
		seriesId,
		geo,
		freqDbNames[strings.ToUpper(freq)],
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id, 0)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}

func (r *FooRepository) GetSeriesSiblingsGeoById(seriesId int64) (geographies []models.DataPortalGeography, err error) {
	//language=MySQL
	rows, err := r.RunQuery(
		`SELECT DISTINCT geographies.fips, geographies.handle, geographies.display_name, geographies.display_name_short
		FROM <%SERIES%> AS series
		JOIN (SELECT name, universe FROM <%SERIES%> where id = ?) AS original_series
		JOIN geographies ON geographies.id = series.geography_id
		WHERE series.universe = original_series.universe
		AND TRIM(TRAILING 'NS' FROM SUBSTRING_INDEX(series.name, '@', 1)) =  /* name prefixes are equal */
			TRIM(TRAILING 'NS' FROM SUBSTRING_INDEX(original_series.name, '@', 1))
		ORDER BY COALESCE(geographies.list_order, 999), geographies.handle`, seriesId)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		geography := models.Geography{}
		err = rows.Scan(
			&geography.FIPS,
			&geography.Handle,
			&geography.Name,
			&geography.ShortName,
		)
		if err != nil {
			continue
		}
		dataPortalGeography := models.DataPortalGeography{Handle: geography.Handle}
		if geography.FIPS.Valid {
			dataPortalGeography.FIPS = geography.FIPS.String
		}
		if geography.Name.Valid {
			dataPortalGeography.Name = geography.Name.String
		}
		if geography.ShortName.Valid {
			dataPortalGeography.ShortName = geography.ShortName.String
		}
		geographies = append(geographies, dataPortalGeography)
	}
	return
}

func (r *FooRepository) GetSeriesSiblingsFreqById(
	seriesId int64,
) (frequencyList []models.DataPortalFrequency, err error) {
	//language=MySQL
	rows, err := r.RunQuery(`SELECT DISTINCT(RIGHT(series.name, 1)) AS freq
	FROM <%SERIES%> AS series
	JOIN (SELECT name, universe FROM <%SERIES%> WHERE id = ?) AS original_series
	WHERE series.universe = original_series.universe
	AND TRIM(TRAILING 'NS' FROM SUBSTRING_INDEX(series.name, '@', 1)) =  /* name prefixes are equal */
	    TRIM(TRAILING 'NS' FROM SUBSTRING_INDEX(original_series.name, '@', 1))
	ORDER BY FIELD(freq, "A", "S", "Q", "M", "W", "D");`, seriesId)
	if err != nil {
		return
	}
	defer rows.Close()
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

func (r *FooRepository) GetSeriesById(seriesId int64, categoryId int64) (dataPortalSeries models.DataPortalSeries, err error) {
	//language=MySQL
	row, err := r.RunQuery(`
		SELECT DISTINCT
		   series_id, series_name, series_universe, series_description, frequency, seasonally_adjusted, seasonal_adjustment,
		   units_long, units_short, data_portal_name, percent, pv.real, source_description, source_link, source_detail_description,
		   table_prefix, table_postfix, measurement_id, measurement_portal_name, NULL, base_year, decimals,
		   geo_fips, geo_handle, geo_display_name, geo_display_name_short
		FROM <%PORTAL%> pv
		WHERE series_id = ? ;`, seriesId) // No MAX() aggregations here because the loop below breaks after first iteration
	if err != nil {
		return
	}
	defer row.Close()
	for row.Next() {
		dataPortalSeries, err = getNextSeriesFromRows(row)
		if err != nil {
			return
		}
		geos, freqs, err := getAllFreqsGeos(r, seriesId, categoryId)
		if err != nil {
			return dataPortalSeries, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		break
	}
	return
}

func (r *FooRepository) GetSeriesByName(name, universe, expand string) (SeriesPkg models.DataPortalSeriesPackage, err error) {
	//language=MySQL
	row, err := r.RunQuery(`
		SELECT
		   s.id, s.name, s.universe, s.description, frequency, seasonally_adjusted, seasonal_adjustment,
		   units.long_label, units.short_label, s.dataPortalName, s.percent, s.real,
		   sources.description, COALESCE(s.source_link, sources.link) AS source_link, source_details.description,
		   NULL, NULL, NULL, NULL, NULL, NULL, decimals,
		   geographies.fips, geographies.handle, geographies.display_name, geographies.display_name_short
		FROM <%SERIES%> s
		JOIN geographies ON geographies.id = s.geography_id
		LEFT JOIN units ON units.id = s.unit_id
		LEFT JOIN sources ON sources.id = s.source_id
		LEFT JOIN source_details ON source_details.id = s.source_detail_id
		WHERE s.name = ? AND s.universe = ? ;`, name, universe)
	if err != nil {
		return
	}
	var series models.DataPortalSeries
	var observations models.SeriesObservations

	defer row.Close()
	if !row.Next() {
		err = row.Err()
		return
	}
	series, err = getNextSeriesFromRows(row)
	if err != nil {
		return
	}
	SeriesPkg.Series = series

	if expand != "" {
		observations, err = r.GetSeriesObservations(SeriesPkg.Series.Id, expand)
		if err != nil {
			return
		}
		SeriesPkg.Observations = &observations
	}
	return
}

// GetSeriesObservations returns an observations struct containing the default transformations.
// It checks the value of percent for the selected series and chooses the appropriate transformations.
//    It has now been turned into a decorator for a more flexible method that allows selection of transformations
func (r *FooRepository) GetSeriesObservations(seriesId int64, expand string) (seriesObservations models.SeriesObservations, err error) {
	return r.GetSeriesTransformations(seriesId, makeBoolSet("all"), expand)
}

func (r *FooRepository) GetSeriesTransformations(seriesId int64, include boolSet, expand string) (seriesObservations models.SeriesObservations, err error) {
	var start, end time.Time
	var percent sql.NullBool
	var universe string
	YOY, YTD, C5MA, MOM := YOYPercentChange, YTDPercentChange, C5MAPercentChange, MOMPercentChange

	//language=MySQL
	err = r.RunQueryRow(`SELECT universe, percent FROM <%SERIES%> WHERE id = ? ;`, seriesId).Scan(&universe, &percent)
	if err != nil {
		return
	}
	if percent.Valid && percent.Bool {
		YOY, YTD, C5MA, MOM = YOYChange, YTDChange, C5MAChange, MOMChange
	}

	var transform models.TransformationResult
	seriesObservations.TransformationResults = make([]models.TransformationResult, 0, 4)

	if include[Levels] || include["all"] {
		transform, err = r.GetTransformation(Levels, seriesId, expand, &start, &end)
		if err != nil {
			return
		}
		seriesObservations.TransformationResults = append(seriesObservations.TransformationResults, transform)
	}
	if include[YOY] || include["all"] && universe != "NTA" {
		transform, err = r.GetTransformation(YOY, seriesId, expand, &start, &end)
		if err != nil {
			return
		}
		seriesObservations.TransformationResults = append(seriesObservations.TransformationResults, transform)
	}
	if include[YTD] || include["all"] && universe != "NTA" {
		transform, err = r.GetTransformation(YTD, seriesId, expand, &start, &end)
		if err != nil {
			return
		}
		seriesObservations.TransformationResults = append(seriesObservations.TransformationResults, transform)
	}
	if include[C5MA] || include["all"] && universe == "NTA" {
		transform, err = r.GetTransformation(C5MA, seriesId, expand, &start, &end)
		if err != nil {
			return
		}
		seriesObservations.TransformationResults = append(seriesObservations.TransformationResults, transform)
	}
	if include[MOM] {
		transform, err = r.GetTransformation(MOM, seriesId, expand, &start, &end)
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
	variadic := make([]interface{}, count)
	for i := range variadic {
		variadic[i] = seriesId
	}
	return variadic
}

func (r *FooRepository) GetTransformation(
	transformation string,
	seriesId int64,
	expand string,
	currentStart *time.Time,
	currentEnd *time.Time,
) (
	transformationResult models.TransformationResult,
	err error,
) {
	var observationStart, observationEnd time.Time
	rows, err := r.RunQuery(transformations[transformation].Statement+" ORDER BY 1;",
		variadicSeriesId(seriesId, transformations[transformation].PlaceholderCount)...,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	obsDates := make([]string, 0, 1000)
	obsValues := make([]string, 0, 1000)
	obsPseudoHist := make([]bool, 0, 1000)

	observation := models.Observation{}
	for rows.Next() {
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
		value := observation.Value.Float64
		if expand == "raw" {
			obsValues = append(obsValues, float64OnlyStringify(value))
		} else {
			obsValues = append(obsValues, float64RoundStringify(value, observation.Decimals))
		}
		obsPseudoHist = append(obsPseudoHist, observation.PseudoHistory.Bool)
	}
	rows.Close()
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

func (r *FooRepository) CreateSeriesPackage(
	id int64,
	universe string,
	categoryId int64,
	forecast string,
	categoryRepository *FooRepository,
) (pkg models.DataPortalSeriesPackage, err error) {

	series, err := r.GetSeriesById(id, categoryId)
	if err != nil {
		return
	}
	pkg.Series = series

	categories, err := categoryRepository.GetAllCategoriesByUniverse(universe)
	if err != nil {
		return
	}
	pkg.Categories = categories

	observations, err := r.GetSeriesObservations(id, "")
	if err != nil {
		return
	}
	pkg.Observations = &observations

	siblings, err := r.GetSeriesSiblingsById(id, forecast, categoryId)
	if err != nil {
		return
	}
	pkg.Siblings = siblings

	forecasts, err := r.GetForecastsById(id)
	if err != nil {
		return
	}
	pkg.Forecasts = forecasts
	return
}

func (r *FooRepository) CreateAnalyzerPackage(
	ids []int64,
	universe string,
	momTransformationOnly bool,
	categoryRepository *FooRepository,
) (pkg models.DataPortalAnalyzerPackage, err error) {

	pkg.InflatedSeries = []models.InflatedSeries{}

	for _, id := range ids {
		series, anErr := r.GetSeriesById(id, 0)
		if anErr != nil {
			err = anErr
			return
		}
		observations, anErr := r.GetSeriesObservations(id, "")
		if momTransformationOnly {
			observations, anErr = r.GetSeriesTransformations(id, makeBoolSet(MOMChange, MOMPercentChange), "")
		}
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

func (r *FooRepository) CreateExportPackage(id int64, expand string) (pkg []models.InflatedSeries, err error) {
	//language=MySQL
	rows, err := r.RunQuery(
		`select s.id, s.universe, s.name, s.dataPortalName from <%SERIES%> s
		 join export_series es on es.series_id = s.id and es.export_id = ?
		 order by es.list_order`, id)
	if err != nil {
		return
	}
	var series models.InflatedSeries
	var dpn sql.NullString
	var observations models.SeriesObservations

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&series.Id, &series.Universe, &series.Name, &dpn)
		if err != nil {
			return
		}
		if dpn.Valid {
			series.Title = dpn.String
		}
		observations, err = r.GetSeriesTransformations(series.Id, makeBoolSet(Levels), expand)
		if err != nil {
			return
		}
		series.Observations = observations
		pkg = append(pkg, series)
	}
	return
}
