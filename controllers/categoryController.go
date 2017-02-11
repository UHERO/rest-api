package controllers

import (
	"encoding/json"
	"net/http"

	"errors"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
	"context"
	"strconv"
	"log"
)

func GetCategory(categoryRepository *data.CategoryRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam, ok := mux.Vars(r)["id"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get category id from request"),
				"Bad request.",
				400,
			)
			return
		}
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}

		category, err := categoryRepository.GetCategoryById(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(CategoryResource{Data: category})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error processing JSON has occurred",
				500,
			)
			return
		}
		ctx := NewContext(context.Background(), j)
		r = r.WithContext(ctx)
		log.Printf("DEBUG: GetCategory: payload is "+string(j))
	}
}

func GetCategories(categoryRepository *data.CategoryRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := categoryRepository.GetAllCategories()
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(CategoriesResource{Data: categories})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error processing JSON has occurred",
				500,
			)
			return
		}
		ctx := NewContext(context.Background(), j)
		if ctx == nil {
			log.Printf("DEBUG: at set point: ctx IS nil")
		} else {
			log.Printf("DEBUG: at set point: ctx NOT nil: %s", ctx)
		}
		*r = *(r.WithContext(ctx))
		//foo := FromContext(r.Context())
		log.Printf("DEBUG: at set point: r is %p", r)
		//log.Printf("DEBUG: GetCategories: payload is "+string(j))
	}
}

func GetCategoryRoots(categoryRepository *data.CategoryRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := categoryRepository.GetCategoryRoots()
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(CategoriesResource{Data: categories})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error processing JSON has occurred",
				500,
			)
			return
		}
		if (string(j) == "foo") {
			log.Printf("bar")
		}
//		rUrl := r.URL.Path+"?"+r.URL.RawQuery
//		context.Set(r, rUrl, j)
//		log.Printf("DEBUG: GetCategoryRoots: rUrl is "+rUrl)
	}
}

func GetCategoriesByName(categoryRepository *data.CategoryRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchText, ok := mux.Vars(r)["searchText"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get searchText from request"),
				"Bad request.",
				400,
			)
			return
		}
		categories, err := categoryRepository.GetCategoriesByName(searchText)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(CategoriesResource{Data: categories})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error processing JSON has occurred",
				500,
			)
			return
		}
		if (string(j) == "foo") {
			log.Printf("bar")
		}
//		rUrl := r.URL.Path+"?"+r.URL.RawQuery
//		context.Set(r, rUrl, j)
	}
}
