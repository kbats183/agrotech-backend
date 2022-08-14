package db

import (
	"database/sql"
	"githab.com/kbats183/argotech/backend/models"
)

func (db Database) GetAllUniversity() (models.UniversityList, error) {
	list := make(models.UniversityList, 0)
	rows, err := db.Conn.Query("SELECT id, name, address, image FROM university ORDER BY id")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var profession models.University
		err := profession.ScanRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, profession)
	}
	return list, nil
}

func (db Database) GetUniversityByID(id int) (models.University, error) {
	profession := models.University{}
	query := `SELECT * FROM university WHERE id = $1;`
	row := db.Conn.QueryRow(query, id)
	switch err := profession.ScanRow(row); err {
	case sql.ErrNoRows:
		return profession, ErrNoMatch
	default:
		return profession, err
	}
}
