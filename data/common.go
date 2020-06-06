package data

import (
	"database/sql"
	"fmt"
	"github.com/UHERO/rest-api/models"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"sort"
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

type FooRepository struct {
	DB				*sql.DB
	PortalView		string
	SeriesView		string
	DataPointView	string
	ReplaceViews	func([]byte) []byte
}

var Rex *regexp.Regexp

func (r *FooRepository) InitializeFoo() *FooRepository {
	var err error
	Rex, err = regexp.Compile(`%[A-Z]+%`)
	if err != nil {
		log.Fatal("Failed to compile the regex")
	}
	r.ReplaceViews = func(str []byte) []byte {
		var result string
		switch string(str) {
		case "%PORTAL%":
			result = r.PortalView
		case "%SERIES%":
			result = r.SeriesView
		case "%DATAPOINTS%":
			result = r.DataPointView
		}
		return []byte(result)
	}
	return r
}

func (r *FooRepository) RunQuery(query string, args ...interface{}) (*sql.Rows, error) {
	query = string(Rex.ReplaceAllFunc([]byte(query), r.ReplaceViews))  // silly that we need to cast back and forth like this :/
	return r.DB.Query(query, args)
}

func (r *FooRepository) RunQueryRow(query string, args ...interface{}) *sql.Row {
	if hasFormat, _ := regexp.MatchString(`%s`, query); hasFormat {
		query = fmt.Sprintf(query, r.PortalView) // plug in portal view name
	}
	return r.DB.QueryRow(query, args)
}

type boolSet map[string]bool

func makeBoolSet(keys ...string) boolSet {
	set := boolSet{}
	for _, key := range keys {
		set[key] = true
	}
	return set
}

func getNextSeriesFromRows(rows *sql.Rows) (dataPortalSeries models.DataPortalSeries, err error) {
	series := models.Series{}
	geography := models.Geography{}
	err = rows.Scan(
		&series.Id,
		&series.Name,
		&series.Universe,
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
		&geography.ShortName,
	)
	if err != nil {
		return
	}
	dataPortalSeries = models.DataPortalSeries{
		Id:       series.Id,
		Name:     series.Name,
		Universe: series.Universe,
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
	if geography.ShortName.Valid {
		dataPortalGeography.ShortName = geography.ShortName.String
	}
	dataPortalSeries.Geography = dataPortalGeography
	return
}

func getAllFreqsGeos(r *FooRepository, seriesId int64, categoryId int64) (
	[]models.DataPortalGeography,
	[]models.DataPortalFrequency,
	error,
) {
	rows, err := r.DB.Query(
		`SELECT DISTINCT 'geo' AS gftype,
			geo.handle AS handle,
			geo.fips AS fips,
			geo.display_name AS display_name,
			geo.display_name_short AS display_name_short,
			geo.list_order AS lorder,
            MIN(pdp.date), MAX(pdp.date)
		FROM measurement_series
		LEFT JOIN measurement_series AS ms ON ms.measurement_id = measurement_series.measurement_id
		LEFT JOIN series_v AS series ON series.id = ms.series_id
		LEFT JOIN geographies geo ON geo.id = series.geography_id
		LEFT JOIN public_data_points pdp on pdp.series_id = series.id
		WHERE pdp.value IS NOT NULL
		AND measurement_series.series_id = ?
		AND NOT (series.restricted OR series.quarantined)
		GROUP BY geo.id, geo.handle, geo.fips, geo.display_name, geo.display_name_short, geo.list_order
		   UNION
		SELECT DISTINCT 'freq' AS gftype,
			RIGHT(series.name, 1) AS handle, null, null, null, null AS lorder, MIN(pdp.date), MAX(pdp.date)
		FROM measurement_series
		LEFT JOIN measurement_series AS ms ON ms.measurement_id = measurement_series.measurement_id
		LEFT JOIN series_v AS series ON series.id = ms.series_id
		LEFT JOIN public_data_points pdp on pdp.series_id = series.id
		WHERE pdp.value IS NOT NULL
		AND measurement_series.series_id = ?
		AND NOT (series.restricted OR series.quarantined)
		GROUP BY RIGHT(series.name, 1)
		ORDER BY gftype, COALESCE(lorder, 999), handle`, seriesId, seriesId)
	if err != nil {
		return nil, nil, err
	}
	geosResult := []models.DataPortalGeography{}
	freqsResult := []models.DataPortalFrequency{}
	for rows.Next() {
		var gftype sql.NullString
		var listOrder sql.NullInt64  // thrown away, but we need to Scan it
		temp := models.Geography{} // Using Geography object as a scan buffer, because it works.
		err = rows.Scan(
			&gftype,
			&temp.Handle,
			&temp.FIPS,
			&temp.Name,
			&temp.ShortName,
			&listOrder,
			&temp.ObservationStart,
			&temp.ObservationEnd,
		)
		if gftype.String == "geo" {
			g := models.DataPortalGeography{Handle: temp.Handle}
			if temp.FIPS.Valid {
				g.FIPS = temp.FIPS.String
			}
			if temp.Name.Valid {
				g.Name = temp.Name.String
			}
			if temp.ShortName.Valid {
				g.ShortName = temp.ShortName.String
			}
			if temp.ObservationStart.Valid {
				g.ObservationStart = &temp.ObservationStart.Time
			}
			if temp.ObservationEnd.Valid {
				g.ObservationEnd = &temp.ObservationEnd.Time
			}
			geosResult = append(geosResult, g)
		} else {
			f := models.DataPortalFrequency{
				Freq:  temp.Handle,
				Label: freqLabel[temp.Handle],
			}
			if temp.ObservationStart.Valid {
				f.ObservationStart = &temp.ObservationStart.Time
			}
			if temp.ObservationEnd.Valid {
				f.ObservationEnd = &temp.ObservationEnd.Time
			}
			freqsResult = append(freqsResult, f)
		}
	}
	sort.Sort(models.ByFrequency(freqsResult))
	return geosResult, freqsResult, err
}

func formatWithYear(formatString string, year int64) string {
	return strings.Replace(formatString, "%Y", strconv.FormatInt(year, 10), -1)
}

func rangeIntersection(start1 time.Time, end1 time.Time, start2 time.Time, end2 time.Time) (iStart *time.Time, iEnd *time.Time) {
	iStart = &start1
	iEnd = &end1
	if !rangesOverlap(start1, end1, start2, end2) {
		return nil, nil
	}
	if start2.After(start1) {
		iStart = &start2
	}
	if end2.Before(end1) {
		iEnd = &end2
	}
	return
}

func rangesOverlap(start1 time.Time, end1 time.Time, start2 time.Time, end2 time.Time) bool {
	return !(end1.Before(start2) || end2.Before(start1))
}
