package data

import (
	"database/sql"
	"github.com/UHERO/rest-api/models"
	"strconv"
	"strings"
)

type CategoryRepository struct {
	DB *sql.DB
}

func (r *CategoryRepository) GetAllCategories() (categories []models.Category, err error) {
	rows, err := r.DB.Query(`
	SELECT cat1.id, MAX(cat1.name), MAX(cat1.ancestry),
	MAX(COALESCE(NULLIF(cat1.default_handle, ''), NULLIF(cat2.default_handle, ''))) AS geo,
	MAX(COALESCE(NULLIF(cat1.default_freq, ''), NULLIF(cat2.default_freq, ''))) AS freq
	FROM categories AS cat1 LEFT JOIN categories AS cat2 ON cat1.ancestry regexp concat('(^|/)', cat2.id, '($|/)') GROUP BY cat1.id;
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
	id, name, ancestry
	FROM categories
	WHERE id = ?;`, id).Scan(
		&category.Id,
		&category.Name,
		&category.Ancestry,
	)
	parentId := getParentId(category.Ancestry)
	return models.Category{
		Id:       category.Id,
		Name:     category.Name,
		ParentId: parentId,
	}, err
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
