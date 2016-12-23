package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"strconv"
	"strings"
	"time"
)

type CategoryRepository struct {
	DB *sql.DB
}

func (r *CategoryRepository) GetAllCategories() (categories []models.Category, err error) {
	rows, err := r.DB.Query(`
SELECT id, name, ancestry, default_handle, default_freq
FROM categories
ORDER BY categories.order;
	`)
	if err != nil {
		return
	}
	for rows.Next() {
		category := models.CategoryWithAncestry{}
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.Ancestry,
			&category.DefaultHandle,
			&category.DefaultFrequency,
		)
		if err != nil {
			return
		}

		parentId := getParentId(category.Ancestry)
		dataPortalCategory := models.Category{
			Id:       category.Id,
			Name:     category.Name,
			ParentId: parentId,
		}
		if category.DefaultHandle.Valid && category.DefaultFrequency.Valid {
			dataPortalCategory.DefaultGeoFreq = &models.GeoFreq{
				Geography: category.DefaultHandle.String,
				Frequency: category.DefaultFrequency.String,
			}
		}
		categories = append(categories, dataPortalCategory)
	}
	return
}

func getParentId(ancestry sql.NullString) (parentId int64) {
	if !ancestry.Valid {
		return
	}
	parents := strings.Split(ancestry.String, "/")
	if len(parents) == 0 {
		return
	}
	parentId, _ = strconv.ParseInt(parents[len(parents)-1], 10, 64)
	return
}

func (r *CategoryRepository) GetCategoryRoots() (categories []models.Category, err error) {
	rows, err := r.DB.Query("SELECT id, name FROM categories WHERE ancestry IS NULL ORDER BY `order`;")
	if err != nil {
		return
	}
	for rows.Next() {
		category := models.Category{}
		err = rows.Scan(
			&category.Id,
			&category.Name,
		)
		if err != nil {
			return
		}
		categories = append(categories, category)
	}
	return
}

func (r *CategoryRepository) GetCategoryById(id int64) (models.Category, error) {
	var category models.CategoryWithAncestry
	err := r.DB.QueryRow(`SELECT
	categories.id, MAX(name), MAX(ancestry),
	MIN(data_points.date) AS start_date, MAX(data_points.date) AS end_date
	FROM categories
	LEFT JOIN data_lists_series ON categories.data_list_id = data_lists_series.data_list_id
 	LEFT JOIN data_points ON data_lists_series.series_id = data_points.series_id
	WHERE categories.id = ? AND data_points.current GROUP BY categories.id;`, id).Scan(
		&category.Id,
		&category.Name,
		&category.Ancestry,
		&category.ObservationStart,
		&category.ObservationEnd,
	)
	parentId := getParentId(category.Ancestry)
	dataPortalCategory := models.Category{
		Id:       category.Id,
		Name:     category.Name,
		ParentId: parentId,
	}
	if category.ObservationStart.Valid && category.ObservationStart.Time.After(time.Time{}) {
		dataPortalCategory.ObservationStart = &category.ObservationStart.Time
	}
	if category.ObservationEnd.Valid && category.ObservationEnd.Time.After(time.Time{}) {
		dataPortalCategory.ObservationEnd = &category.ObservationEnd.Time
	}
	return dataPortalCategory, err
}

func (r *CategoryRepository) GetCategoriesByName(name string) (categories []models.Category, err error) {
	fuzzyString := "%" + name + "%"
	rows, err := r.DB.Query("SELECT id, name, ancestry FROM categories WHERE LOWER(name) LIKE ? ORDER BY `order`;", fuzzyString)
	if err != nil {
		return
	}
	for rows.Next() {
		category := models.CategoryWithAncestry{}
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.Ancestry,
		)
		if err != nil {
			return
		}
		parentId := getParentId(category.Ancestry)
		categories = append(categories, models.Category{
			Id:       category.Id,
			Name:     category.Name,
			ParentId: parentId,
		})
	}
	return
}

func (r *CategoryRepository) GetChildrenOf(id int64) (categories []models.Category, err error) {
	rows, err := r.DB.Query("SELECT id, name, parent_id FROM categories WHERE parent_id = ? ORDER BY `order;`", id)
	if err != nil {
		return
	}
	for rows.Next() {
		category := models.Category{}
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.ParentId,
		)
		if err != nil {
			return
		}
		categories = append(categories, category)
	}
	return
}
