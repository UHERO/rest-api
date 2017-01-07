package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"errors"
)

var freqLabel map[string]string = map[string]string{
	"A": "Annual",
	"S": "Semiannual",
	"Q": "Quarterly",
	"M": "Monthly",
	"W": "Weekly",
	"D": "Daily",
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
		&series.UnitsLabel,
		&series.UnitsLabelShort,
		&series.DataPortalName,
		&series.Percent,
		&series.Real,
		&geography.FIPS,
		&geography.Handle,
		&geography.Name,
	)
	if err != nil {
		return
	}
	dataPortalSeries = models.DataPortalSeries{
		Id:             series.Id,
		Name:           series.Name,
		FrequencyShort: series.Name[len(series.Name)-1:],
	}
	dataPortalSeries.Frequency = freqLabel[dataPortalSeries.FrequencyShort]
	if series.DataPortalName.Valid {
		dataPortalSeries.Title = series.DataPortalName.String
	}
	if series.Description.Valid {
		dataPortalSeries.Description = series.Description.String
	}
	if series.SeasonallyAdjusted.Valid {
		dataPortalSeries.SeasonallyAdjusted = &series.SeasonallyAdjusted.Bool
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
		&series.UnitsLabel,
		&series.UnitsLabelShort,
		&series.DataPortalName,
		&series.Percent,
		&series.Real,
		&geography.FIPS,
		&geography.Handle,
		&geography.Name,
	)
	if err != nil {
		return dataPortalSeries, errors.New("Series restricted or does not exist.")
	}
	dataPortalSeries = models.DataPortalSeries{
		Id:             series.Id,
		Name:           series.Name,
		FrequencyShort: series.Name[len(series.Name)-1:],
	}
	dataPortalSeries.Frequency = freqLabel[dataPortalSeries.FrequencyShort]
	if series.DataPortalName.Valid {
		dataPortalSeries.Title = series.DataPortalName.String
	}
	if series.Description.Valid {
		dataPortalSeries.Description = series.Description.String
	}
	if series.SeasonallyAdjusted.Valid {
		dataPortalSeries.SeasonallyAdjusted = &series.SeasonallyAdjusted.Bool
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
