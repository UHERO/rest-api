package data

import (
	"database/sql"
	"log"
	"time"
	"crypto/rand"
	"github.com/uhero/rest-api/models"
	"encoding/base64"
)

type ApplicationRepository struct {
	DB *sql.DB
}

func (r *ApplicationRepository) Create(userName string, application *models.Application) (numRows int64, err error) {
	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil {
		panic(err)
	}
	application.Key = base64.URLEncoding.EncodeToString(rb)
	stmt, err := r.DB.Prepare(`INSERT INTO api_applications(name, hostname, 'key', github_nickname, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?);`)
	if err != nil {
		panic(err)
	}
	res, err := stmt.Exec(
		application.Name,
		application.Hostname,
		application.Key,
		userName,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		panic(err)
	}
	application.Id, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}
	numRows, err = res.RowsAffected()
	return
}

func (r *ApplicationRepository) Update(userName string, application *models.Application) (numRows int64, err error) {
	stmt, err := r.DB.Prepare(`UPDATE api_applications SET
	name = ?,
	hostname = ?,
	'key' = ?,
	updated_at = ?,
	WHERE id = ? and github_nickname = ? LIMIT 1;`)
	if err != nil {
		panic(err)
	}
	res, err := stmt.Exec(
		application.Name,
		application.Hostname,
		application.Key,
		time.Now(),
		application.Id,
		userName,
	)
	if err != nil {
		panic(err)
	}
	numRows, err = res.RowsAffected()
	log.Printf("%d Rows updated", numRows)
	return
}

func (r *ApplicationRepository) Delete(userName string, id int64) (numRows int64, err error) {
	stmt, err := r.DB.Prepare(`DELETE FROM api_applications WHERE id = ? and github_nickname = ?;`)
	if err != nil {
		panic(err)
	}
	res, err := stmt.Exec(id, userName)
	if err != nil {
		panic(err)
	}
	numRows, err = res.RowsAffected()
	log.Printf("%d Rows deleted", numRows)
	return
}

func (r *ApplicationRepository) GetAll(userName string) (applications []models.Application, err error) {
	rows, err := r.DB.Query(`SELECT
	id, name, hostname, 'key'
	FROM api_applications WHERE github_nickname = ?;`, userName)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		application := models.Application{}
		err = rows.Scan(
			&application.Id,
			&application.Name,
			&application.Hostname,
			&application.Key,
		)
		if err != nil {
			panic(err)
		}
		applications = append(applications, application)
	}
	return
}

func (r *ApplicationRepository) GetById(userName string, id int64) (application models.Application, err error) {
	err = r.DB.QueryRow(`SELECT
	id, name, hostname, 'key'
	FROM api_applications
	WHERE id = ? AND github_nickname = ?;`, id, userName).Scan(
		&application.Id,
		&application.Name,
		&application.Hostname,
		&application.Key,
	)
	return
}
