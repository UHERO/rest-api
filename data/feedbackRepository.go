package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
)

type FeedbackRepository struct {
	DB *sql.DB
}

func (r *FeedbackRepository) CreateFeedback(feedback *models.Feedback) (err error) {
	return
}

