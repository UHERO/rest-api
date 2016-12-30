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
	Id               int64      `json:"id"`
	Name             string     `json:"name"`
	ParentId         int64      `json:"parentId,omitempty"`
	DefaultGeoFreq   *GeoFreq   `json:"defaults,omitempty"`
	ObservationStart *time.Time `json:"observationStart,omitempty"`
	ObservationEnd   *time.Time `json:"observationEnd,omitempty"`
}

type GeoFreq struct {
	Geography string `json:"geo,omitempty"`
	Frequency string `json:"freq,omitempty"`
}

type CategoryWithAncestry struct {
	Id               int64
	Name             string
	Ancestry         sql.NullString
	DefaultHandle    sql.NullString
	DefaultFrequency sql.NullString
	ObservationStart NullTime
	ObservationEnd   NullTime
}

type SearchSummary struct {
	SearchText       string     `json:"q"`
	DefaultGeoFreq   *GeoFreq   `json:"defaults,omitempty"`
	GeoFreqs         map[string][]string  `json:"geoFreqs"`
	FreqGeos         map[string][]string  `json:"freqGeos"`
	ObservationStart *time.Time `json:"observationStart,omitempty"`
	ObservationEnd   *time.Time `json:"observationStart,omitempty"`
}

type Geography struct {
	FIPS   sql.NullString `json:"fips"`
	Name   sql.NullString `json:"name"`
	Handle string         `json:"handle"`
}

type DataPortalGeography struct {
	FIPS   string `json:"fips"`
	Name   string `json:"name"`
	Handle string `json:"handle"`
}

type FrequencyResult struct {
	Freq  string `json:"freq"`
	Label string `json:"label"`
}

type Frequency struct {
	Freq  string
	Label sql.NullString
}

type Series struct {
	Id                 int64
	Name               string
	Description        sql.NullString
	Frequency          sql.NullString
	SeasonallyAdjusted sql.NullBool
	UnitsLabel         sql.NullString
	UnitsLabelShort    sql.NullString
	DataPortalName     sql.NullString
	Percent            sql.NullBool
	Real               sql.NullBool
}

type DataPortalSeries struct {
	Id                 int64               `json:"id"`
	Name               string              `json:"name"`
	Title              string              `json:"title,omitempty"`
	Description        string              `json:"description,omitempty"`
	Frequency          string              `json:"frequency,omitempty"`
	FrequencyShort     string              `json:"frequencyShort,omitempty"`
	SeasonallyAdjusted *bool               `json:"seasonallyAdjusted,omitempty"`
	UnitsLabel         string              `json:"unitsLabel,omitEmpty"`
	UnitsLabelShort    string              `json:"unitsLabelShort,omitEmpty"`
	Geography          DataPortalGeography `json:"geography,omitEmpty"`
	Percent            *bool               `json:"percent,omitempty"`
	Real               *bool               `json:"real,omitempty"`
}

type InflatedSeries struct {
	DataPortalSeries
	Observations SeriesObservations `json:"seriesObservations"`
}

type Observation struct {
	Date  time.Time
	Value sql.NullFloat64
	PseudoHistory bool
}

type DataPortalObservation struct {
	Date  time.Time
	Value float64
	PseudoHistory bool
}

type SeriesObservations struct {
	ObservationStart      time.Time              `json:"observationStart"`
	ObservationEnd        time.Time              `json:"observationEnd"`
	OrderBy               string                 `json:"orderBy"`
	SortOrder             string                 `json:"sortOrder"`
	TransformationResults []TransformationResult `json:"transformationResults"`
}

type TransformationResult struct {
	Transformation string                  `json:"transformation"`
	Observations   []DataPortalObservation `json:"observations"`
}

func (o *DataPortalObservation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Date  string `json:"date"`
		Value string `json:"value"`
		PseudoHistory bool `json:"pseudoHistory,omitempty"`
	}{
		Date:  formatDate(o.Date),
		Value: fmt.Sprintf("%.4f", o.Value),
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
