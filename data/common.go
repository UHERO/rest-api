package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"sort"
	"strconv"
	"strings"
)

var freqLabel map[string]string = map[string]string{
	"A": "Annual",
	"S": "Semiannual",
	"Q": "Quarterly",
	"M": "Monthly",
	"W": "Weekly",
	"D": "Daily",
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
		&series.UnitsLabel,
		&series.UnitsLabelShort,
		&series.DataPortalName,
		&series.Percent,
		&series.Real,
		&series.SourceDescription,
		&series.SourceLink,
		&series.SourceDetails,
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
	if series.SeasonallyAdjusted.Valid && series.Name[len(series.Name)-1:] != "A" {
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
		&series.UnitsLabel,
		&series.UnitsLabelShort,
		&series.DataPortalName,
		&series.Percent,
		&series.Real,
		&series.SourceDescription,
		&series.SourceLink,
		&series.SourceDetails,
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
	if series.SeasonallyAdjusted.Valid && series.Name[len(series.Name)-1:] != "A" {
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

func getFreqGeoCombinations(r *SeriesRepository, seriesId int64) (
	[]models.GeographyFrequencies,
	[]models.FrequencyGeographies,
	error,
) {
	rows, err := r.DB.Query(`SELECT geographies.fips, geographies.display_name_short, geofreq.geo, geofreq.freq
	FROM (SELECT MAX(SUBSTRING_INDEX(SUBSTR(name, LOCATE('@', name) + 1), '.', 1)) as geo,
		   MAX(RIGHT(name, 1)) as freq
	FROM (SELECT series.name
		FROM measurement_series
		LEFT JOIN measurement_series AS ms ON ms.measurement_id = measurement_series.measurement_id
		LEFT JOIN series ON series.id = ms.series_id
		WHERE measurement_series.series_id = ?) AS s
	GROUP BY SUBSTR(name, LOCATE('@', name) + 1) ORDER BY COUNT(*) DESC) as geofreq
	LEFT JOIN geographies ON geographies.handle = geofreq.geo;`, seriesId)
	if err != nil {
		return nil, nil, err
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

	return geoFreqsResult, freqGeosResult, err
}

func formatWithYear(formatString string, year int64) string {
	return strings.Replace(formatString, "%Y", strconv.FormatInt(year, 10), -1)
}
