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

	// GET - /package/category
	CategoryPackage struct {
		Data models.DataPortalCategoryPackage `json:"data"`
	}

	// GET - /category
	CategoriesResource struct {
		Data []models.Category `json:"data"`
	}

	// GET - /category/measurements
	MeasurementListResource struct {
		Data []models.Measurement `json:"data"`
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
		Data []models.DataPortalFrequency `json:"data"`
	}

	ForecastListResource struct {
		Data []string `json:"data"`
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

	// GET - /search
	SearchSummaryResource struct {
		Data models.SearchSummary `json:"data"`
	}

	// POST - /feedback
	FeedbackResource struct {
		Data models.Feedback `json:"data"`
	}

	// GET - /package/series
	SeriesPackage struct {
		Data models.DataPortalSeriesPackage `json:"data"`
	}

	// GET - /package/search
	SearchPackage struct {
		Data models.DataPortalSearchPackage `json:"data"`
	}

	// GET - /package/analyzer
	AnalyzerPackage struct {
		Data models.DataPortalAnalyzerPackage `json:"data"`
	}

	// GET - /package/export
	ExportPackage struct {
		Data []models.InflatedSeries `json:"data"`
	}
)
