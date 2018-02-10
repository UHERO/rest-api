package controllers

import (
	"encoding/json"
	"net/http"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/UHERO/rest-api/models"
	"github.com/gorilla/mux"
	"errors"
)

func GetSeriesByGroupId(
	seriesRepository *data.SeriesRepository,
	cacheRepository *data.CacheRepository,
	groupType data.GroupType,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByGroup(id, groupType)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetInflatedSeriesByGroupId(
	seriesRepository *data.SeriesRepository,
	cacheRepository *data.CacheRepository,
	groupType data.GroupType,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetInflatedSeriesByGroup(id, groupType)
		returnInflatedSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesByGroupIdAndGeoHandle(
	seriesRepository *data.SeriesRepository,
	cacheRepository *data.CacheRepository,
	groupType data.GroupType,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, ok := getIdAndGeo(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByGroupAndGeo(id, geoHandle, groupType)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesByGroupIdAndFreq(
	seriesRepository *data.SeriesRepository,
	cacheRepository *data.CacheRepository,
	groupType data.GroupType,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, freq, ok := getIdAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByGroupAndFreq(id, freq, groupType)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesByGroupIdGeoHandleAndFreq(
	seriesRepository *data.SeriesRepository,
	cacheRepository *data.CacheRepository,
	groupType data.GroupType,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesByGroupGeoAndFreq(id, geoHandle, freq, groupType)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetInflatedSeriesByGroupIdGeoAndFreq(
	seriesRepository *data.SeriesRepository,
	cacheRepository *data.CacheRepository,
	groupType data.GroupType,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetInflatedSeriesByGroupGeoAndFreq(id, geoHandle, freq, groupType)
		returnInflatedSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesById(seriesRepository *data.SeriesRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		series, err := seriesRepository.GetSeriesById(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(SeriesResource{Data: series})
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
		WriteCache(r, cacheRepository, j)
	}
}

func GetSeriesSiblingsById(seriesRepository *data.SeriesRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsById(id)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesSiblingsByIdAndFreq(seriesRepository *data.SeriesRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, freq, ok := getIdAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdAndFreq(id, freq)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesSiblingsByIdAndGeo(seriesRepository *data.SeriesRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geo, ok := getIdAndGeo(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdAndGeo(id, geo)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesSiblingsByIdGeoAndFreq(seriesRepository *data.SeriesRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geo, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdGeoAndFreq(id, geo, freq)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesSiblingsFreqById(seriesRepository *data.SeriesRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		frequencyList, err := seriesRepository.GetSeriesSiblingsFreqById(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(FrequencyListResource{Data: frequencyList})
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
		WriteCache(r, cacheRepository, j)
	}
}

func GetFreqByCategoryId(seriesRepository *data.SeriesRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		frequencyList, err := seriesRepository.GetFreqByCategory(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(FrequencyListResource{Data: frequencyList})
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
		WriteCache(r, cacheRepository, j)
	}
}

func GetSeriesObservations(seriesRepository *data.SeriesRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		series, err := seriesRepository.GetSeriesObservations(id)
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
			return
		}
		j, err := json.Marshal(ObservationList{Data: series})
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
		WriteCache(r, cacheRepository, j)
	}
}

func GetSeriesPackage(
	seriesRepository *data.SeriesRepository,
	categoryRepository *data.CategoryRepository,
	cacheRepository *data.CacheRepository,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pkg := models.DataPortalSeriesPackage{}
		id, ok := getId(w, r)
		if !ok {
			return
		}
		series, err := seriesRepository.GetSeriesById(id)
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
		pkg.Series = series

		categories, err := categoryRepository.GetAllCategoriesByUniverse(series.Universe)
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
		pkg.Categories = categories

		observations, err := seriesRepository.GetSeriesObservations(id)
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
		pkg.Observations = observations

		siblings, err := seriesRepository.GetSeriesSiblingsById(id)
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
		pkg.Siblings = siblings

		j, err := json.Marshal(SeriesPackage{Data: pkg})
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error processing JSON has occurred", 500)
			return
		}
		WriteResponse(w, j)
		WriteCache(r, cacheRepository, j)
	}
}

func GetAnalyzerPackage(
	categoryRepository *data.CategoryRepository,
	seriesRepository *data.SeriesRepository,
	cacheRepository *data.CacheRepository,
) func(http.ResponseWriter, *http.Request) {
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
		ids, ok := getIdsList(w, r)
		if !ok {
			common.DisplayAppError(
				w,
				errors.New("Couldn't get id list from request"),
				"Bad request.",
				400,
			)
			ok = false
			return
		}
		pkg, err := categoryRepository.CreateAnalyzerPackage(universe, ids, seriesRepository)
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}

		j, err := json.Marshal(AnalyzerPackage{Data: pkg})
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error processing JSON has occurred", 500)
			return
		}
		WriteResponse(w, j)
		WriteCache(r, cacheRepository, j)
	}
}
