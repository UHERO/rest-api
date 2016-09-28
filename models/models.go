package models

import "database/sql"

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
