package controllers

import (
	"github.com/uhero/rest-api/models"
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
)
