package models

import (
	"database/sql"
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
	Id                 string `json:"id"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	Frequency          string `json:"frequency,omitempty"`
	SeasonallyAdjusted bool   `json:"seasonallyAdjusted,omitempty"`
	UnitsLabel         string `json:"unitsLabel,omitEmpty"`
	UnitsLabelShort    string `json:"unitsLabelShort,omitEmpty"`
}
