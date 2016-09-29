package data

import (
	"database/sql"
	"github.com/uhero/rest-api/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestGetAllCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	category1 := models.CategoryWithAncestry{
		Id:       1,
		Name:     "Summary",
		Ancestry: sql.NullString{Valid: false},
	}
	category2 := models.CategoryWithAncestry{
		Id:       2,
		Name:     "Income",
		Ancestry: sql.NullString{Valid: true, String: "1"},
	}
	categoryResult := sqlmock.NewRows([]string{"id", "name", "ancestry"}).
		AddRow(category1.Id, category1.Name, nil).
		AddRow(category2.Id, category2.Name, category2.Ancestry.String)
	mock.ExpectQuery("SELECT id, name, ancestry FROM categories").
		WillReturnRows(categoryResult)

	categoryRepository := CategoryRepository{DB: db}

	categories, err := categoryRepository.GetAllCategories()
	if err != nil {
		t.Fail()
	}
	if len(categories) != 2 ||
		categories[0].Name != category1.Name ||
		categories[1].Id != category2.Id {
		t.Fail()
	}
}
