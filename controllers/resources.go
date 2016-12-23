package controllers

import (
	"github.com/UHERO/rest-api/models"
)

type (
	// for GET - /applications
	ApplicationsResource struct {
		Data []models.Application `json:"data"`
	}

	// for POST/PUT - /applications
	// for GET - /applications/id
	ApplicationResource struct {
		Data models.Application `json:"data"`
	}

	// for main display
	UserResource struct {
		User         string
		Applications []models.Application
	}
)

// for core API, all prefixed by /v1
type (
	// GET - /category/id
	CategoryResource struct {
		Data models.Category `json:"data"`
	}

	// GET - /category
	CategoriesResource struct {
		Data []models.Category `json:"data"`
	}

	// GET - /category/series
	SeriesListResource struct {
		Data []models.DataPortalSeries `json:"data"`
	}
	InflatedSeriesListResource struct {
		Data []models.InflatedSeries `json:"data"`
	}

	// GET - /series/siblings/freq
	FrequencyListResource struct {
		Data []models.FrequencyResult `json:"data"`
	}

	// GET - /series
	SeriesResource struct {
		Data models.DataPortalSeries `json:"data"`
	}

	// GET - /series/observations
	ObservationList struct {
		Data models.SeriesObservations `json:"data"`
	}

	// GET - /geo
	GeographiesResource struct {
		Data []models.DataPortalGeography `json:"data"`
	}
)
