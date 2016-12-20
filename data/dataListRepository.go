package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"log"
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

	// find requested measurements

	// create data_list_measurements entries
	return
}

// Loop over each id to find all the siblings
func (r *DataListRepository) siblingIdsFromSeriesIds(seriesIds []int64) (siblings map[int]bool, err error) {
	//for seriesId := range seriesIds {
	//	// runs a query for each entry in the list
	//	//rows, err := r.DB.Query(siblingIds, seriesId)
	//	//if err != nil {
	//	//	return
	//	//}
	//	//for rows.Next() {
	//	//	var id int
	//	//	err = rows.Scan(&id)
	//	//	if err != nil {
	//	//		return
	//	//	}
	//	//	siblings[id] = true
	//	//}
	//}
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
