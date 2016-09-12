package models

type Application struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	APIKey      string `json:"apiKey"`
}
