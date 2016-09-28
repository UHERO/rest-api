package data

import (
	"github.com/uhero/rest-api/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
)

func TestGetAllCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	category1 := models.Category{
		Id:       1,
		Name:     "Summary",
		ParentId: 0,
	}
	category2 := models.Category{
		Id:       2,
		Name:     "Income",
		ParentId: 1,
	}
	categoryResult := sqlmock.NewRows([]string{"id", "name", "parentId"}).
		AddRow(category1.Id, category1.Name, category1.ParentId).
		AddRow(category2.Id, category2.Name, category2.ParentId)
	mock.ExpectQuery("SELECT id, name, parent_id FROM categories").
		WillReturnRows(categoryResult)

	categoryRepository := CategoryRepository{DB: db}
	
	categories, err := categoryRepository.GetAll()
	if err != nil {
		t.Fail()
	}
	if len(categories) != 2 ||
		!reflect.DeepEqual(categories[0], category1) ||
		!reflect.DeepEqual(categories[1], category2) {
		t.Fail()
	}
}

