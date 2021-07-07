package controllers

import (
	"encoding/json"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"net/http"
	"strings"
)

func GetSeriesByGroupId(
	seriesRepository *data.FooRepository,
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
	seriesRepository *data.FooRepository,
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
	seriesRepository *data.FooRepository,
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
	seriesRepository *data.FooRepository,
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
	seriesRepository *data.FooRepository,
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
	seriesRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
	groupType data.GroupType,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geoHandle, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		forecast, _ := getStrParam(r, "forecast")
		if !ok {
			forecast = "@"  // a regex that will match any series name
		}

		seriesList, err := seriesRepository.GetInflatedSeriesByGroupGeoAndFreq(id, geoHandle, freq, forecast, groupType)
		returnInflatedSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesById(seriesRepository *data.FooRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		series, err := seriesRepository.GetSeriesById(id, 0)
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


func GetSeriesByName(seriesRepository *data.FooRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := getStrParam(r, "name")
		if !ok {
			return
		}
		name = strings.Replace(name, "-", "&", -1)  // replace all "-" placeholders with "&"
		universe, ok := getStrParam(r, "universe")
		if !ok {
			universe = "UHERO"
		}
		startDate, _ := getStrParam(r, "start_from")
		expand, _ := getStrParam(r, "exp")

		seriesPkg, err := seriesRepository.GetSeriesByName(name, universe, startDate, expand == "true")
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
		j, err := json.Marshal(SeriesPackage{Data: seriesPkg})
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error processing JSON has occurred", 500)
			return
		}
		WriteResponse(w, j)
		WriteCache(r, cacheRepository, j)
	}
}

func GetSeriesSiblingsById(seriesRepository *data.FooRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		catId, ok := getIntParam(r, "cat")
		if !ok {
			catId = 0
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsById(id, catId)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesSiblingsByIdAndFreq(seriesRepository *data.FooRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, freq, ok := getIdAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdAndFreq(id, freq)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesSiblingsByIdAndGeo(seriesRepository *data.FooRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geo, ok := getIdAndGeo(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdAndGeo(id, geo)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesSiblingsByIdGeoAndFreq(seriesRepository *data.FooRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, geo, freq, ok := getIdGeoAndFreq(w, r)
		if !ok {
			return
		}
		seriesList, err := seriesRepository.GetSeriesSiblingsByIdGeoAndFreq(id, geo, freq)
		returnSeriesList(seriesList, err, w, r, cacheRepository)
	}
}

func GetSeriesSiblingsFreqById(seriesRepository *data.FooRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
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

func GetFreqByCategoryId(seriesRepository *data.FooRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		forecast, ok := getStrParam(r, "forecast")
		if !ok {
			forecast = "@"  // a regex that will match any series name
		}
		frequencyList, err := seriesRepository.GetFreqByCategory(id, forecast)
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

func GetForecastByCategoryId(seriesRepository *data.FooRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		forecastList, err := seriesRepository.GetForecastByCategory(id)
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
		j, err := json.Marshal(ForecastListResource{Data: forecastList})
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error processing JSON has occurred",  500)
			return
		}
		WriteResponse(w, j)
		WriteCache(r, cacheRepository, j)
	}
}

func GetSeriesObservations(seriesRepository *data.FooRepository, cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		series, err := seriesRepository.GetSeriesObservations(id, "")
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
	seriesRepository *data.FooRepository,
	categoryRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		id, ok := getId(w, r)
		if !ok {
			return
		}
		catId, ok := getIntParam(r, "cat")
		if !ok {
			catId = 0
		}
		universe, ok := getStrParam(r, "universe_text")
		if !ok {
			return
		}
		startDate, _ := getStrParam(r, "start_from")

		pkg, err := seriesRepository.CreateSeriesPackage(id, universe, catId, startDate, categoryRepository)
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
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
	categoryRepository *data.FooRepository,
	seriesRepository *data.FooRepository,
	cacheRepository *data.CacheRepository,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ids, ok := getIdsList(w, r)
		if !ok {
			return
		}
		universe, ok := getStrParam(r, "universe_text")
		if !ok {
			return
		}
		pkg, err := seriesRepository.CreateAnalyzerPackage(ids, universe, categoryRepository)
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

func GetExportPackage(seriesRepo *data.FooRepository, c *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := getId(w, r)
		if !ok {
			return
		}
		pkg, err := seriesRepo.CreateExportPackage(id)
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error has occurred", 500)
			return
		}
		j, err := json.Marshal(ExportPackage{Data: pkg})
		if err != nil {
			common.DisplayAppError(w, err, "An unexpected error processing JSON has occurred", 500)
			return
		}
		WriteResponse(w, j)
		WriteCache(r, c, j)
	}
}
