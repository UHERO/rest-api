package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"log"
	"strings"
	"time"
)

type DataListRepository struct {
	DB *sql.DB
}

func (r *DataListRepository) CreateDataList(userId int64, dataList *models.DataList) (numRows int64, err error) {
	stmt, err := r.DB.Prepare(`INSERT INTO
	data_list(name, created_by, owned_by, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?);`)
	if err != nil {
		return
	}
	res, err := stmt.Exec(
		dataList.Name,
		userId,
		userId,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return
	}
	dataList.Id, err = res.LastInsertId()
	if err != nil {
		return
	}
	numRows, err = res.RowsAffected()
	if err != nil {
		return
	}
	siblings, err := r.siblingIdsFromSeriesIds(dataList.SeriesIds)
	if err != nil {
		return
	}

	// loop over each id
	siblingSql := "INSERT INTO data_lists_series(data_lists_id, series_id) VALUES "
	vals := make([]int64, 5)
	for siblingId := range siblings {
		siblingSql += "(?, ?),"
		vals = append(vals, dataList.Id, siblingId)
	}
	strings.TrimSuffix(siblingSql, ",")
	stmt, err = r.DB.Prepare(siblingSql)
	if err != nil {
		return
	}
	res, err = stmt.Exec(vals...)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// Loop over each id to find all the siblings
func (r *DataListRepository) siblingIdsFromSeriesIds(seriesIds []int64) (siblings map[int]bool, err error) {
	for seriesId := range seriesIds {
		// runs a query for each entry in the list
		rows, err := r.DB.Query(siblingIds, seriesId)
		if err != nil {
			return
		}
		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			if err != nil {
				return
			}
			siblings[id] = true
		}
	}
	return
}

func (r *DataListRepository) UpdateDataList(userId int64, dataList *models.DataList) (numRows int64, err error) {
	log.Printf("Update DataList %d for user %d", dataList.Id, userId)
	// update the data_lists entry
	// if no rows were affected then return
	// drop data_lists_series entries for this data_list
	// create data_lists_series entries
	return
}

func (r *DataListRepository) DeleteDataList(userId int64, dataListId int64) (numRows int64, err error) {
	return
}

func (r *DataListRepository) GetAllDataLists(userId int64) (dataLists []models.DataList, err error) {
	return
}
