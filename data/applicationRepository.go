package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"log"
	"time"
)

type AppRepository interface {
	CreateApplication(string, *models.Application) (int64, error)
	UpdateApplication(string, *models.Application) (int64, error)
	DeleteApplication(string, int64) (int64, error)
	GetAllApplications(string) ([]models.Application, error)
}

type ApplicationRepository struct {
	DB *sql.DB
}

func (r *FooRepository) CreateApplication(username string, application *models.Application) (numRows int64, err error) {
	application.APIKey, err = CreateNewApiKey(32)
	if err != nil {
		return
	}
	stmt, err := r.DB.Prepare(`INSERT INTO api_applications(name, hostname, api_key, github_nickname, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return
	}
	t := time.Now()
	res, err := stmt.Exec(
		application.Name,
		application.Hostname,
		application.APIKey,
		username,
		t,
		t,
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

func (r *FooRepository) UpdateApplication(username string, application *models.Application) (numRows int64, err error) {
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

func (r *FooRepository) DeleteApplication(username string, id int64) (numRows int64, err error) {
	stmt, err := r.DB.Prepare(`DELETE FROM api_applications WHERE id = ? and github_nickname = ?;`)
	if err != nil {
		return
	}
	res, err := stmt.Exec(id, username)
	if err != nil {
		return
	}
	numRows, err = res.RowsAffected()
	log.Printf("%d Rows deleted", numRows)
	return
}

func (r *FooRepository) GetAllApplications(username string) (applications []models.Application, err error) {
	rows, err := r.DB.Query(`SELECT id, name, hostname, api_key FROM api_applications WHERE github_nickname = ?;`, username)
	if err != nil {
		return
	}
	defer rows.Close()
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

func (r *FooRepository) GetApplicationsByApiKey(apiKey string) (applications []models.Application, err error) {
	rows, err := r.DB.Query(`SELECT id, name, hostname, api_key FROM api_applications WHERE api_key = ?;`, apiKey)
	if err != nil {
		return
	}
	defer rows.Close()
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

func (r *FooRepository) GetApplicationById(username string, id int64) (application models.Application, err error) {
	err = r.DB.QueryRow(`SELECT id, name, hostname, api_key FROM api_applications WHERE id = ? AND github_nickname = ?;`,
		id, username).Scan(
			&application.Id,
			&application.Name,
			&application.Hostname,
			&application.APIKey,
	)
	return
}
