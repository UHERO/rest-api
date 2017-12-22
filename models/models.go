package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Application struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	APIKey   string `json:"apiKey"`
}

type Category struct {
	Id                   int64                   `json:"id"`
	Name                 string                  `json:"name"`
	Universe             string                  `json:"universe"`
	ParentId             int64                   `json:"parentId,omitempty"`
	Defaults	     *CategoryDefaults	     `json:"defaults,omitempty"`
	Current		     *CurrentGeoFreq	     `json:"current,omitempty"`
	Geographies          *[]DataPortalGeography  `json:"geos,omitempty"`
	Frequencies          *[]DataPortalFrequency  `json:"freqs,omitempty"`
	ObservationStart     *time.Time              `json:"observationStart,omitempty"`
	ObservationEnd       *time.Time              `json:"observationEnd,omitempty"`
}

type CategoryDefaults struct {
	Geography            *DataPortalGeography    `json:"geo,omitempty"`
	Frequency            *DataPortalFrequency    `json:"freq,omitempty"`
	ObservationStart     *time.Time              `json:"observationStart,omitempty"`
	ObservationEnd       *time.Time              `json:"observationEnd,omitempty"`
}

type CurrentGeoFreq struct {
	Geography            string		`json:"geo,omitempty"`
	Frequency            string		`json:"freq,omitempty"`
	ObservationStart     *time.Time		`json:"observationStart,omitempty"`
	ObservationEnd       *time.Time		`json:"observationEnd,omitempty"`
}

type CategoryWithAncestry struct {
	Id			int64
	Name			string
	Universe		string
	Ancestry		sql.NullString
	ParentId		sql.NullInt64
	DefaultGeoHandle	sql.NullString
	DefaultGeoFIPS		sql.NullString
	DefaultGeoName		sql.NullString
	DefaultGeoShortName	sql.NullString
	DefaultFrequency	sql.NullString
	ObservationStart	NullTime
	ObservationEnd  	NullTime
}

type DataPortalCategoryPackage struct {
	CatSubTree	[]CategoryWithInflatedSeries	`json:"categories"`
	NavCategories	[]Category			`json:"navCategories"`
}

type CategoryWithInflatedSeries struct {
	Category
	Series []InflatedSeries		`json:"series"`
}

type SearchSummary struct {
	SearchText           string                  `json:"q"`
	DefaultGeo           *DataPortalGeography    `json:"defaultGeo,omitempty"`
	DefaultFreq          *DataPortalFrequency    `json:"defaultFreq,omitempty"`
	Geographies          *[]DataPortalGeography  `json:"geos,omitempty"`
	Frequencies          *[]DataPortalFrequency  `json:"freqs,omitempty"`
	ObservationStart     *time.Time              `json:"observationStart"`
	ObservationEnd       *time.Time              `json:"observationEnd"`
}

type Geography struct {
	FIPS             sql.NullString `json:"fips"`
	Name             sql.NullString `json:"name"`
	ShortName        sql.NullString `json:"shortName"`
	Handle           string         `json:"handle"`
	ObservationStart NullTime
	ObservationEnd   NullTime
}

type Frequency struct {
	Freq  string
	Label sql.NullString
	ObservationStart NullTime
	ObservationEnd   NullTime
}

type DataPortalGeography struct {
	FIPS             string     `json:"fips,omitempty"`
	Name             string     `json:"name,omitempty"`
	ShortName        string     `json:"shortName,omitempty"`
	Handle           string     `json:"handle"`
	ObservationStart *time.Time `json:"observationStart,omitempty"`
	ObservationEnd   *time.Time `json:"observationEnd,omitempty"`
}

type DataPortalFrequency struct {
	Freq             string     `json:"freq"`
	Label            string     `json:"label,omitempty"`
	ObservationStart *time.Time `json:"observationStart,omitempty"`
	ObservationEnd   *time.Time `json:"observationEnd,omitempty"`
}

// ByGeography/ByFrequency implement sort.Interface
type ByGeography []DataPortalGeography
type ByFrequency []DataPortalFrequency
type stringSlice []string

func (s stringSlice) indexOf(stringToFind string) int {
	for key, value := range s {
		if value == stringToFind {
			return key
		}
	}
	return -1
}

var FreqOrder = stringSlice{"A", "S", "Q", "M", "W", "D"}

func (a ByFrequency) Len() int      { return len(a) }
func (a ByFrequency) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByFrequency) Less(i, j int) bool {
	return FreqOrder.indexOf(a[i].Freq) < FreqOrder.indexOf(a[j].Freq)
}

func (a ByGeography) Len() int      { return len(a) }
func (a ByGeography) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByGeography) Less(i, j int) bool {
	return a[i].Handle < a[j].Handle
}

type Series struct {
	Id                 int64
	Name               string
	Universe	   string
	Description        sql.NullString
	Frequency          sql.NullString
	SeasonallyAdjusted sql.NullBool
	SeasonalAdjustment sql.NullString
	UnitsLabel         sql.NullString
	UnitsLabelShort    sql.NullString
	DataPortalName     sql.NullString
	Percent            sql.NullBool
	Real               sql.NullBool
	Decimals           sql.NullInt64
	BaseYear           sql.NullInt64
	SourceDescription  sql.NullString
	SourceLink         sql.NullString
	SourceDetails      sql.NullString
	Indent             sql.NullString
	TablePrefix        sql.NullString
	TablePostfix       sql.NullString
	MeasurementId      sql.NullInt64
	MeasurementName    sql.NullString
}

type Measurement struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Indent int    `json:"indent,omitempty"`
}

type DataPortalSeriesPackage struct {
	Series		DataPortalSeries	`json:"series"`
	Categories	[]Category		`json:"categories"`
	Observations	SeriesObservations	`json:"observations"`
	Siblings	[]DataPortalSeries	`json:"siblings"`
}

type DataPortalSeries struct {
	Id                          int64                   `json:"id"`
	Name                        string                  `json:"name"`
	Universe		    string                  `json:"universe"`
	Title                       string                  `json:"title,omitempty"`
	Description                 string                  `json:"description,omitempty"`
	MeasurementId               int64                   `json:"measurementId,omitempty"`
	MeasurementName             string                  `json:"measurementName,omitempty"`
	Frequency                   string                  `json:"frequency,omitempty"`
	FrequencyShort              string                  `json:"frequencyShort,omitempty"`
	SeasonallyAdjusted          *bool                   `json:"seasonallyAdjusted,omitempty"`
	SeasonalAdjustment          string                  `json:"seasonalAdjustment,omitempty"`
	UnitsLabel                  string                  `json:"unitsLabel,omitEmpty"`
	UnitsLabelShort             string                  `json:"unitsLabelShort,omitEmpty"`
	Geography                   DataPortalGeography     `json:"geography,omitEmpty"`
	Percent                     *bool                   `json:"percent,omitempty"`
	Real                        *bool                   `json:"real,omitempty"`
	BaseYear                    *int64                  `json:"baseYear,omitempty"`
	BaseYearDeprecated          *int64                  `json:"base_year,omitempty"`
	Decimals                    *int64                  `json:"decimals,omitempty"`
	SourceDescription           string                  `json:"sourceDescription,omitempty"`
	SourceLink                  string                  `json:"sourceLink,omitempty"`
	SourceDescriptionDeprecated string                  `json:"source_description,omitempty"`
	SourceLinkDeprecated        string                  `json:"source_link,omitempty"`
	SourceDetails               string                  `json:"sourceDetails,omitempty"`
	Indent                      int                     `json:"indent,omitempty"`
	TablePrefix                 string                  `json:"tablePrefix"`
	TablePostfix                string                  `json:"tablePostfix"`
	Geographies                 *[]DataPortalGeography  `json:"geos,omitempty"`
	Frequencies		    *[]DataPortalFrequency  `json:"freqs,omitempty"`
}

type InflatedSeries struct {
	DataPortalSeries
	Observations SeriesObservations `json:"seriesObservations"`
}

type Observation struct {
	Date          time.Time
	Value         sql.NullFloat64
	PseudoHistory sql.NullBool
	Decimals      int
}

type DataPortalObservation struct {
	Date          time.Time
	Value         float64
	PseudoHistory *bool
}

type SeriesObservations struct {
	ObservationStart      time.Time              `json:"observationStart"`
	ObservationEnd        time.Time              `json:"observationEnd"`
	OrderBy               string                 `json:"orderBy"`
	SortOrder             string                 `json:"sortOrder"`
	TransformationResults []TransformationResult `json:"transformationResults"`
}

type TransformationResult struct {
	Transformation string        `json:"transformation"`
	ObservationDates   []string  `json:"dates"`
	ObservationValues  []string  `json:"values"`
	ObservationPHist   []bool    `json:"pseudoHistory"`
}

type Feedback struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Feedback string `json:"feedback"`
}

func (o *DataPortalObservation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Date          string `json:"date"`
		Value         string `json:"value"`
		PseudoHistory *bool  `json:"pseudoHistory,omitempty"`
	}{
		Date:          formatDate(o.Date),
		Value:         fmt.Sprintf("%.4f", o.Value),
		PseudoHistory: o.PseudoHistory,
	})
}

// solid idea from stack overflow: http://choly.ca/post/go-json-marshalling/
func (so SeriesObservations) MarshalJSON() ([]byte, error) {
	type Alias SeriesObservations
	return json.Marshal(&struct {
		ObservationStart string `json:"observationStart"`
		ObservationEnd   string `json:"observationEnd"`
		Alias
	}{
		ObservationStart: formatDate(so.ObservationStart),
		ObservationEnd:   formatDate(so.ObservationEnd),
		Alias:            (Alias)(so),
	})
}

func formatDate(date time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", date.Year(), date.Month(), date.Day())
}

// from https://github.com/lib/pq/blob/master/encode.go
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	if nt.Time.IsZero() {
		nt.Valid = false
	}
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
