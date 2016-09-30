package models

import (
	"database/sql"
	"time"
)

type Application struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	APIKey   string `json:"apiKey"`
}

type Category struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	ParentId int64  `json:"parentId,omitempty"`
}

type CategoryWithAncestry struct {
	Id       int64
	Name     string
	Ancestry sql.NullString
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
}

type DataPortalSeries struct {
	Id                 int64 `json:"id"`
	Name               string `json:"name"`
	Title              string `json:"title,omitempty"`
	Description        string `json:"description,omitempty"`
	Frequency          string `json:"frequency,omitempty"`
	SeasonallyAdjusted bool   `json:"seasonallyAdjusted,omitempty"`
	UnitsLabel         string `json:"unitsLabel,omitEmpty"`
	UnitsLabelShort    string `json:"unitsLabelShort,omitEmpty"`
}

type Observation struct {
	Date time.Time `json:"date"`
	Value float64 `json:"value"`
}

type SeriesObservations struct {
	ObservationStart time.Time `json:"observationStart"`
	ObservationEnd time.Time `json:"observationEnd"`
	OrderBy string `json:"orderBy"`
	SortOrder string `json:"sortOrder"`
	TransformationResults []TransformationResult `json:"transformationResults"`
}

type TransformationResult struct {
	Transformation string `json:"transformation"`
	Observations []Observation `json:"observations"`

}