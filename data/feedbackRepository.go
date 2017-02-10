package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
)

type FeedbackRepository struct {
	DB *sql.DB
}

func (r *FeedbackRepository) CreateFeedback(application *models.Application) (err error) {
	return
}

