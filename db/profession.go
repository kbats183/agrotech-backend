package db

import (
	"database/sql"
	"githab.com/kbats183/argotech/backend/models"
)

func (db Database) GetAllProfession() (models.ProfessionList, error) {
	list := make(models.ProfessionList, 0)
	rows, err := db.Conn.Query("SELECT * FROM professions ORDER BY id")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var profession models.Profession
		err := profession.ScanRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, profession)
	}
	return list, nil
}

func (db Database) GetProfessionByID(id int) (models.Profession, error) {
	profession := models.Profession{}
	query := `SELECT * FROM professions WHERE id = $1;`
	row := db.Conn.QueryRow(query, id)
	switch err := profession.ScanRow(row); err {
	case sql.ErrNoRows:
		return profession, ErrNoMatch
	default:
		return profession, err
	}
}

func (db Database) GetAllProfessionWithRating(userAuth models.UserAuth) (models.ProfessionWithRatingList, error) {
	list := make(models.ProfessionWithRatingList, 0)
	rows, err := db.Conn.Query("SELECT p.*, fp.id IS NOT NULL as \"is_favourite\" FROM professions p LEFT JOIN favourite_professions fp on p.id = fp.profession_id LEFT JOIN users u on u.id = fp.user_id WHERE u.login = $1 OR u.login is NULL ORDER BY p.id;", userAuth)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var profession models.ProfessionWithRating
		err := profession.ScanRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, profession)
	}
	return list, nil
}

func (db Database) GetProfessionWithRatingByID(userAuth models.UserAuth, id int) (models.ProfessionWithRating, error) {
	profession := models.ProfessionWithRating{}
	query := `SELECT p.*, fp.id IS NOT NULL as "is_favourite" FROM professions p LEFT JOIN favourite_professions fp on p.id = fp.profession_id LEFT JOIN users u on u.id = fp.user_id WHERE (u.login = $1 OR u.login is NULL) AND p.id = $2;`
	row := db.Conn.QueryRow(query, userAuth, id)
	switch err := profession.ScanRow(row); err {
	case sql.ErrNoRows:
		return profession, ErrNoMatch
	default:
		return profession, err
	}
}

func (db Database) AddProfessionFavourite(userID models.UserID, professionID int) error {
	rows := db.Conn.QueryRow("INSERT INTO favourite_professions(profession_id, user_id) VALUES ($2, $1) ON CONFLICT DO NOTHING;", userID, professionID)
	return rows.Err()
}

func (db Database) DeleteProfessionFavourite(userID models.UserID, professionID int) error {
	rows := db.Conn.QueryRow("DELETE FROM favourite_professions WHERE user_id=$1 AND profession_id=$2;", userID, professionID)
	return rows.Err()
}
