package data

import (
	"github.com/UHERO/rest-api/models"
)

func (r *FooRepository) GetAllForecasts() (forecastList models.ForecastList, err error) {
	//language=MySQL
	rows, err := r.RunQuery(`SELECT DISTINCT SUBSTRING_INDEX(SUBSTRING_INDEX(series_name, '@', 1), '&', -1) AS fc
							 FROM <%PORTAL%> pv
							 WHERE pv.series_universe = 'FC';`)
	if err != nil {
		return
	}
	defer rows.Close()
	var fcName, freqFull string
	for rows.Next() {
		err = rows.Scan(&fcName)
		if err != nil {
			return
		}
		code := freqDbCodes[freqFull]
		forecastList = append(forecastList, fcName)
	}
	return
}
