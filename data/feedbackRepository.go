package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
)

type FeedbackRepository struct {
	DB *sql.DB
}

func (r *FeedbackRepository) CreateFeedback(feedback models.Feedback) (numRows int64, err error) {
	stmt, err := r.DB.Prepare(`INSERT INTO user_feedbacks(name, email, feedback, created_at, updated_at)
	VALUES (?, ?, ?, NOW(), NOW());`)
	if err != nil {
		return
	}
	res, err := stmt.Exec(feedback.Name, feedback.Email, feedback.Feedback)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
