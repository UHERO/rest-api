package data

import (
	"database/sql"
	"log"
	"time"
	"crypto/rand"
	"github.com/uhero/rest-api/models"
	"encoding/base64"
)

type Repository interface {
	Create(string, *models.Application) (int64, error)
	Update(string, *models.Application) (int64, error)
	Delete(string, int64) (int64, error)
	GetAll(string) ([]models.Application, error)
}

type ApplicationRepository struct {
	DB *sql.DB
}

func (r *ApplicationRepository) Create(userName string, application *models.Application) (numRows int64, err error) {
	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil {
		return
	}
	application.APIKey = base64.URLEncoding.EncodeToString(rb)
	stmt, err := r.DB.Prepare(`INSERT INTO api_applications(name, hostname, api_key, github_nickname, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return
	}
	res, err := stmt.Exec(
		application.Name,
		application.Hostname,
		application.APIKey,
		userName,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return
	}
	application.Id, err = res.LastInsertId()
	if err != nil {
		return
	}
	numRows, err = res.RowsAffected()
	return
}

func (r *ApplicationRepository) Update(username string, application *models.Application) (numRows int64, err error) {
	stmt, err := r.DB.Prepare(`UPDATE api_applications SET
	name = COALESCE(?, name),
	hostname = COALESCE(?, hostname),
	updated_at = ?
	WHERE id = ? and github_nickname = ? LIMIT 1;`)
	if err != nil {
		return
	}
	log.Printf(
		"name = %s\nhostname = %s\napi_key = %s\nid = %d\nusername = %s\n",
		application.Name,
		application.Hostname,
		application.APIKey,
		application.Id,
		username,
	)
	res, err := stmt.Exec(
		application.Name,
		application.Hostname,
		time.Now(),
		application.Id,
		username,
	)
	if err != nil {
		return
	}
	numRows, err = res.RowsAffected()
	log.Printf("%d Rows updated", numRows)
	return
}

func (r *ApplicationRepository) Delete(userName string, id int64) (numRows int64, err error) {
	stmt, err := r.DB.Prepare(`DELETE FROM api_applications WHERE id = ? and github_nickname = ?;`)
	if err != nil {
		return
	}
	res, err := stmt.Exec(id, userName)
	if err != nil {
		return
	}
	numRows, err = res.RowsAffected()
	log.Printf("%d Rows deleted", numRows)
	return
}

func (r *ApplicationRepository) GetAll(userName string) (applications []models.Application, err error) {
	rows, err := r.DB.Query(`SELECT
	id, name, hostname, api_key
	FROM api_applications WHERE github_nickname = ?;`, userName)
	if err != nil {
		return
	}
	for rows.Next() {
		application := models.Application{}
		err = rows.Scan(
			&application.Id,
			&application.Name,
			&application.Hostname,
			&application.APIKey,
		)
		if err != nil {
			return
		}
		applications = append(applications, application)
	}
	return
}

func (r *ApplicationRepository) GetById(userName string, id int64) (application models.Application, err error) {
	err = r.DB.QueryRow(`SELECT
	id, name, hostname, api_key
	FROM api_applications
	WHERE id = ? AND github_nickname = ?;`, id, userName).Scan(
		&application.Id,
		&application.Name,
		&application.Hostname,
		&application.APIKey,
	)
	return
}
