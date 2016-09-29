package data

import (
	"database/sql"
	"github.com/uhero/rest-api/models"
)

type SeriesRepository struct {
	DB *sql.DB
}

func (r *SeriesRepository) GetSeriesByCategory(categoryId int64) (seriesList []models.DataPortalSeries, err error) {
	rows, err := r.DB.Query(`SELECT id, name, description, frequency,
	seasonally_adjusted, unitsLabel, unitsLabelShort, dataPortalName
	FROM series WHERE
	(SELECT list FROM data_lists JOIN categories WHERE categories.data_list_id = data_lists.id AND categories.id = ?)
	LIKE CONCAT('%', left(name, locate("@", name)), '%');`, categoryId)
	if err != nil {
		return
	}
	for rows.Next() {
		series := models.Series{}
		err = rows.Scan(
			&series.Id,
			&series.Name,
			&series.Description,
			&series.Frequency,
			&series.SeasonallyAdjusted,
			&series.UnitsLabel,
			&series.UnitsLabelShort,
			&series.DataPortalName,
		)
		if err != nil {
			return
		}
		dataPortalSeries := models.DataPortalSeries{Id: series.Name}
		if series.DataPortalName.Valid {
			dataPortalSeries.Name = series.DataPortalName.String
		}
		if series.Description.Valid {
			dataPortalSeries.Description = series.Description.String
		}
		if series.Frequency.Valid {
			dataPortalSeries.Frequency = series.Frequency.String
		}
		if series.SeasonallyAdjusted.Valid {
			dataPortalSeries.SeasonallyAdjusted = series.SeasonallyAdjusted.Bool
		}
		if series.UnitsLabel.Valid {
			dataPortalSeries.UnitsLabel = series.UnitsLabel.String
		}
		if series.UnitsLabelShort.Valid {
			dataPortalSeries.UnitsLabelShort = series.UnitsLabelShort.String
		}
		seriesList = append(seriesList, dataPortalSeries)
	}
	return
}
