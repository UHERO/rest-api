package controllers

import (
	"encoding/json"
	"net/http"

	"errors"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
	"strconv"
)

func GetCategory(categoryRepository *data.CategoryRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
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
		WriteResponse(w, j)
		WriteCache(r, c, j)
	}
}

func GetCategoryByIdGeoFreq(categoryRepository *data.CategoryRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
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
		geo, ok := mux.Vars(r)["geo"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get geography from request"),
				"Bad request.",
				400,
			)
			return
		}
		freq, ok := mux.Vars(r)["freq"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get frequency from request"),
				"Bad request.",
				400,
			)
			return
		}
		category, err := categoryRepository.GetCategoryByIdGeoFreq(id, geo, freq)
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
		WriteResponse(w, j)
		WriteCache(r, c, j)
	}
}

func GetCategories(categoryRepository *data.CategoryRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
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
		WriteResponse(w, j)
		WriteCache(r, c, j)
	}
}

func GetCategoryRoots(categoryRepository *data.CategoryRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
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
		WriteResponse(w, j)
		WriteCache(r, c, j)
	}
}

func GetCategoriesByName(categoryRepository *data.CategoryRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
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
		WriteResponse(w, j)
		WriteCache(r, c, j)
	}
}

func GetCategoriesByUniverse(categoryRepository *data.CategoryRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		universe, ok := mux.Vars(r)["universe_text"]
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get universe handle from request"),
				"Bad request.",
				400,
			)
			ok = false
			return
		}
		categories, err := categoryRepository.GetAllCategoriesByUniverse(universe)
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
		WriteResponse(w, j)
		WriteCache(r, c, j)
	}
}
