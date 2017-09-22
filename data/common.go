package data

import (
	"database/sql"
	"sort"
	"strconv"
	"strings"

	"github.com/UHERO/rest-api/models"
)

var freqLabel map[string]string = map[string]string{
	"A": "Annual",
	"S": "Semiannual",
	"Q": "Quarterly",
	"M": "Monthly",
	"W": "Weekly",
	"D": "Daily",
}

var freqDbNames map[string]string = map[string]string{
	"A": "year",
	"S": "semi",
	"Q": "quarter",
	"M": "month",
	"W": "week",
	"D": "day",
}

var indentationLevel map[string]int = map[string]int{
	"indent0": 0,
	"indent1": 1,
	"indent2": 2,
	"indent3": 3,
}

func getNextSeriesFromRows(rows *sql.Rows) (dataPortalSeries models.DataPortalSeries, err error) {
	series := models.Series{}
	geography := models.Geography{}
	err = rows.Scan(
		&series.Id,
		&series.Name,
		&series.Description,
		&series.Frequency,
		&series.SeasonallyAdjusted,
		&series.SeasonalAdjustment,
		&series.UnitsLabel,
		&series.UnitsLabelShort,
		&series.DataPortalName,
		&series.Percent,
		&series.Real,
		&series.SourceDescription,
		&series.SourceLink,
		&series.SourceDetails,
		&series.TablePrefix,
		&series.TablePostfix,
		&series.MeasurementId,
		&series.MeasurementName,
		&series.Indent,
		&series.BaseYear,
		&series.Decimals,
		&geography.FIPS,
		&geography.Handle,
		&geography.Name,
	)
	if err != nil {
		return
	}
	dataPortalSeries = models.DataPortalSeries{
		Id:   series.Id,
		Name: series.Name,
	}
	dataPortalSeries.FrequencyShort = dataPortalSeries.Name[len(dataPortalSeries.Name)-1:]
	dataPortalSeries.Frequency = freqLabel[dataPortalSeries.FrequencyShort]
	if series.DataPortalName.Valid {
		dataPortalSeries.Title = series.DataPortalName.String
	}
	if series.Description.Valid {
		dataPortalSeries.Description = series.Description.String
	}
	if series.SeasonallyAdjusted.Valid && dataPortalSeries.FrequencyShort != "A" {
		dataPortalSeries.SeasonallyAdjusted = &series.SeasonallyAdjusted.Bool
	}
	if series.SeasonalAdjustment.Valid {
		dataPortalSeries.SeasonalAdjustment = series.SeasonalAdjustment.String
	}
	if series.UnitsLabel.Valid {
		dataPortalSeries.UnitsLabel = series.UnitsLabel.String
	}
	if series.UnitsLabelShort.Valid {
		dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
	}
	if series.Percent.Valid {
		dataPortalSeries.Percent = &series.Percent.Bool
	}
	if series.Real.Valid {
		dataPortalSeries.Real = &series.Real.Bool
	}
	if series.SourceDescription.Valid {
		dataPortalSeries.SourceDescription = series.SourceDescription.String
		dataPortalSeries.SourceDescriptionDeprecated = series.SourceDescription.String
	}
	if series.SourceLink.Valid {
		dataPortalSeries.SourceLink = series.SourceLink.String
		dataPortalSeries.SourceLinkDeprecated = series.SourceLink.String
	}
	if series.SourceDetails.Valid {
		dataPortalSeries.SourceDetails = series.SourceDetails.String
	}
	if series.TablePrefix.Valid {
		dataPortalSeries.TablePrefix = series.TablePrefix.String
	}
	if series.TablePostfix.Valid {
		dataPortalSeries.TablePostfix = series.TablePostfix.String
	}
	if series.MeasurementId.Valid {
		dataPortalSeries.MeasurementId = series.MeasurementId.Int64
	}
	if series.MeasurementName.Valid {
		dataPortalSeries.MeasurementName = series.MeasurementName.String
	}
	if series.Decimals.Valid {
		dataPortalSeries.Decimals = &series.Decimals.Int64
	}
	if series.BaseYear.Valid {
		dataPortalSeries.Title = formatWithYear(dataPortalSeries.Title, series.BaseYear.Int64)
		dataPortalSeries.Description = formatWithYear(dataPortalSeries.Description, series.BaseYear.Int64)
		dataPortalSeries.UnitsLabel = formatWithYear(dataPortalSeries.UnitsLabel, series.BaseYear.Int64)
		dataPortalSeries.UnitsLabelShort = formatWithYear(dataPortalSeries.UnitsLabelShort, series.BaseYear.Int64)
		dataPortalSeries.BaseYear = &series.BaseYear.Int64
		dataPortalSeries.BaseYearDeprecated = &series.BaseYear.Int64
	}
	if series.Indent.Valid {
		dataPortalSeries.Indent = indentationLevel[series.Indent.String]
	}
	dataPortalGeography := models.DataPortalGeography{Handle: geography.Handle}
	if geography.FIPS.Valid {
		dataPortalGeography.FIPS = geography.FIPS.String
	}
	if geography.Name.Valid {
		dataPortalGeography.Name = geography.Name.String
	}
	dataPortalSeries.Geography = dataPortalGeography
	return
}

func getNextSeriesFromRow(row *sql.Row) (dataPortalSeries models.DataPortalSeries, err error) {
	series := models.Series{}
	geography := models.Geography{}
	err = row.Scan(
		&series.Id,
		&series.Name,
		&series.Description,
		&series.Frequency,
		&series.SeasonallyAdjusted,
		&series.SeasonalAdjustment,
		&series.UnitsLabel,
		&series.UnitsLabelShort,
		&series.DataPortalName,
		&series.Percent,
		&series.Real,
		&series.SourceDescription,
		&series.SourceLink,
		&series.SourceDetails,
		&series.TablePrefix,
		&series.TablePostfix,
		&series.MeasurementId,
		&series.MeasurementName,
		&series.BaseYear,
		&series.Decimals,
		&geography.FIPS,
		&geography.Handle,
		&geography.Name,
	)
	if err != nil {
		return dataPortalSeries, err
	}
	dataPortalSeries = models.DataPortalSeries{
		Id:   series.Id,
		Name: series.Name,
	}
	dataPortalSeries.FrequencyShort = dataPortalSeries.Name[len(dataPortalSeries.Name)-1:]
	dataPortalSeries.Frequency = freqLabel[dataPortalSeries.FrequencyShort]
	if series.DataPortalName.Valid {
		dataPortalSeries.Title = series.DataPortalName.String
	}
	if series.Description.Valid {
		dataPortalSeries.Description = series.Description.String
	}
	if series.SeasonallyAdjusted.Valid && dataPortalSeries.FrequencyShort != "A" {
		dataPortalSeries.SeasonallyAdjusted = &series.SeasonallyAdjusted.Bool
	}
	if series.SeasonalAdjustment.Valid {
		dataPortalSeries.SeasonalAdjustment = series.SeasonalAdjustment.String
	}
	if series.UnitsLabel.Valid {
		dataPortalSeries.UnitsLabel = series.UnitsLabel.String
	}
	if series.UnitsLabelShort.Valid {
		dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
	}
	if series.Percent.Valid {
		dataPortalSeries.Percent = &series.Percent.Bool
	}
	if series.Real.Valid {
		dataPortalSeries.Real = &series.Real.Bool
	}
	if series.SourceDescription.Valid {
		dataPortalSeries.SourceDescription = series.SourceDescription.String
		dataPortalSeries.SourceDescriptionDeprecated = series.SourceDescription.String
	}
	if series.SourceLink.Valid {
		dataPortalSeries.SourceLink = series.SourceLink.String
		dataPortalSeries.SourceLinkDeprecated = series.SourceLink.String
	}
	if series.SourceDetails.Valid {
		dataPortalSeries.SourceDetails = series.SourceDetails.String
	}
	if series.TablePrefix.Valid {
		dataPortalSeries.TablePrefix = series.TablePrefix.String
	}
	if series.TablePostfix.Valid {
		dataPortalSeries.TablePostfix = series.TablePostfix.String
	}
	if series.MeasurementId.Valid {
		dataPortalSeries.MeasurementId = series.MeasurementId.Int64
	}
	if series.MeasurementName.Valid {
		dataPortalSeries.MeasurementName = series.MeasurementName.String
	}
	if series.Decimals.Valid {
		dataPortalSeries.Decimals = &series.Decimals.Int64
	}
	if series.BaseYear.Valid && series.BaseYear.Int64 > 0 {
		dataPortalSeries.Title = formatWithYear(dataPortalSeries.Title, series.BaseYear.Int64)
		dataPortalSeries.Description = formatWithYear(dataPortalSeries.Description, series.BaseYear.Int64)
		dataPortalSeries.UnitsLabel = formatWithYear(dataPortalSeries.UnitsLabel, series.BaseYear.Int64)
		dataPortalSeries.UnitsLabelShort = formatWithYear(dataPortalSeries.UnitsLabelShort, series.BaseYear.Int64)
		dataPortalSeries.BaseYear = &series.BaseYear.Int64
		dataPortalSeries.BaseYearDeprecated = &series.BaseYear.Int64
	}
	dataPortalGeography := models.DataPortalGeography{Handle: geography.Handle}
	if geography.FIPS.Valid {
		dataPortalGeography.FIPS = geography.FIPS.String
	}
	if geography.Name.Valid {
		dataPortalGeography.Name = geography.Name.String
	}
	dataPortalSeries.Geography = dataPortalGeography
	return
}

func getAllFreqsGeos(r *SeriesRepository, seriesId int64) (
	[]models.DataPortalGeography,
	[]models.DataPortalFrequency,
	error,
) {
	rows, err := r.DB.Query(
		`SELECT DISTINCT 'geo' AS gftype, geo.handle, geo.fips, geo.display_name_short
		FROM measurement_series
		LEFT JOIN measurement_series AS ms ON ms.measurement_id = measurement_series.measurement_id
		LEFT JOIN series ON series.id = ms.series_id
		LEFT JOIN geographies geo on geo.id = series.geography_id
		LEFT JOIN public_data_points pdp on pdp.series_id = series.id
		WHERE pdp.value IS NOT NULL
		AND measurement_series.series_id = ?
			UNION
		SELECT DISTINCT 'freq' AS gftype, RIGHT(series.name, 1) AS handle, null, null
		FROM measurement_series
		LEFT JOIN measurement_series AS ms ON ms.measurement_id = measurement_series.measurement_id
		LEFT JOIN series ON series.id = ms.series_id
		LEFT JOIN public_data_points pdp on pdp.series_id = series.id
		WHERE pdp.value IS NOT NULL
		AND measurement_series.series_id = ?
		ORDER BY 1,2 ;`, seriesId, seriesId)
	if err != nil {
		return nil, nil, err
	}
	geosResult := []models.DataPortalGeography{}
	freqsResult := []models.DataPortalFrequency{}
	for rows.Next() {
		var gftype sql.NullString
		err = rows.Scan(&gftype)

		if gftype == "geo" {
			geography := models.DataPortalGeography{}
			geo_temp := models.Geography{}
			err = rows.Scan(
				&gftype,
				&geography.Handle,
				&geo_temp.FIPS,
				&geo_temp.Name,
			)
			if geo_temp.FIPS.Valid {
				geography.FIPS = geo_temp.FIPS.String
			}
			if geo_temp.Name.Valid {
				geography.Name = geo_temp.Name.String
			}
			geosResult = append(geosResult, geography)
		} else {
			frequency := models.DataPortalFrequency{}
			err = rows.Scan(
				&gftype,
				&frequency.Freq,
			)
			frequency.Label = freqLabel[frequency.Freq]
			freqsResult = append(freqsResult, frequency)
		}
	}
	return geosResult, freqsResult, err
}

func formatWithYear(formatString string, year int64) string {
	return strings.Replace(formatString, "%Y", strconv.FormatInt(year, 10), -1)
}
