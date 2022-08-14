package db

import (
	"fmt"
	"githab.com/kbats183/argotech/backend/models"
)

func (db Database) GetAllVacancies() (models.VacancyList, error) {
	list := make(models.VacancyList, 0)
	rows, err := db.Conn.Query("SELECT * FROM vacancies ORDER BY id")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var profession models.Vacancy
		types, err := rows.ColumnTypes()
		if err != nil {
			return nil, err
		}
		fmt.Println(types)
		err = profession.ScanRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, profession)
	}
	return list, nil
}
