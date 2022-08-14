package db

import (
	"database/sql"
	"githab.com/kbats183/argotech/backend/models"
)

func (db Database) GetAllProfession() (models.ProfessionShortInfoList, error) {
	list := make(models.ProfessionShortInfoList, 0)
	rows, err := db.Conn.Query("SELECT id, name, description, short_description, image FROM professions ORDER BY id")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var profession models.ProfessionShortInfo
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

func (db Database) GetAllProfessionWithRating(userAuth models.UserAuth) (models.ProfessionShortInfoWithRatingList, error) {
	list := make(models.ProfessionShortInfoWithRatingList, 0)
	query := `SELECT p.id,
       p.name,
       p.description,
       p.short_description,
       p.image,
       EXISTS(SELECT fp.profession_id
              FROM favourite_professions fp
                       INNER JOIN users u on u.id = fp.user_id
              WHERE u.login = $1
                AND fp.profession_id = p.id) as "is_favourite"
FROM professions p
ORDER BY p.id;`
	rows, err := db.Conn.Query(query, userAuth)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var profession models.ProfessionShortInfoWithRating
		err := profession.ScanRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, profession)
	}
	return list, nil
}

func (db Database) GetAllFavouriteProfession(userAuth models.UserAuth) (models.ProfessionShortInfoWithRatingList, error) {
	list := make(models.ProfessionShortInfoWithRatingList, 0)
	rows, err := db.Conn.Query("SELECT p.id, p.name, p.description, p.short_description, p.image, true as \"is_favourite\" FROM professions p INNER JOIN favourite_professions fp on p.id = fp.profession_id LEFT JOIN users u on u.id = fp.user_id WHERE u.login = $1 ORDER BY p.id;", userAuth)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var profession models.ProfessionShortInfoWithRating
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
	query := `SELECT p.*, exists(SELECT fp.id FROM favourite_professions fp INNER JOIN users u on u.id = fp.user_id WHERE p.id = fp.profession_id AND u.login = $1) as "is_favourite" FROM professions p WHERE p.id = $2;`
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
