package data

import (
	"github.com/UHERO/rest-api/models"
)

func (r *FooRepository) GetAllForecasts() (forecastList models.ForecastList, err error) {
	//language=MySQL
	rows, err := r.RunQuery(`SELECT DISTINCT SUBSTRING_INDEX(SUBSTRING_INDEX(name, '@', 1), '&', -1) AS fc
							 FROM <%SERIES%>
							 WHERE universe = 'FC';`)
	if err != nil {
		return
	}
	defer rows.Close()
	var fcName string
	for rows.Next() {
		err = rows.Scan(&fcName)
		if err != nil {
			return
		}
		forecastList = append(forecastList, fcName)
	}
	return
}

func (r *FooRepository) GetAllPortalForecasts() (forecastList models.ForecastList, err error) {
	//language=MySQL
	rows, err := r.RunQuery(`SELECT DISTINCT SUBSTRING_INDEX(SUBSTRING_INDEX(series_name, '@', 1), '&', -1) AS fc
							 FROM <%PORTAL%>
							 WHERE series_universe = 'FC';`)
	if err != nil {
		return
	}
	defer rows.Close()
	var fcName string
	for rows.Next() {
		err = rows.Scan(&fcName)
		if err != nil {
			return
		}
		forecastList = append(forecastList, fcName)
	}
	return
}

func (r *FooRepository) GetForecastSeries(forecast, freq string) (seriesList []models.InflatedSeries, err error) {
	rows, err := r.RunQuery(`
		SELECT 
`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		dataPortalSeries, scanErr := getNextSeriesFromRows(rows)
		if scanErr != nil {
			return seriesList, scanErr
		}
		geos, freqs, err := getAllFreqsGeos(r, dataPortalSeries.Id, 0)
		if err != nil {
			return seriesList, err
		}
		dataPortalSeries.Geographies = &geos
		dataPortalSeries.Frequencies = &freqs
		seriesObservations, scanErr := r.GetSeriesObservations(dataPortalSeries.Id, "")
		if scanErr != nil {
			return seriesList, scanErr
		}
		inflatedSeries := models.InflatedSeries{dataPortalSeries, seriesObservations}
		seriesList = append(seriesList, inflatedSeries)
	}
	return
}
