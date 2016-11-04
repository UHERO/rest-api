package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
)

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
		&geography.FIPS,
		&geography.Handle,
		&geography.Name,
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
		&geography.FIPS,
		&geography.Handle,
		&geography.Name,
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
