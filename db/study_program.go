package db

import (
	"database/sql"
	"githab.com/kbats183/argotech/backend/models"
)

func (db Database) GetAllProgramsByProfessionID(professionID int, userID models.UserID) (models.StudyProgramsShortList, error) {
	list := make(models.StudyProgramsShortList, 0)
	query := `SELECT prog.id,
       spec.id,
       spec.name,
       university_id,
       u.name,
       u.image,
       exams,
       exists(SELECT fsp.id
              FROM favourite_study_programs fsp
              WHERE prog.id = fsp.study_programs_id AND fsp.user_id = $2) as "is_favourite"
FROM study_programs prog
         INNER JOIN specialization spec on spec.id = prog.specialization_id
         INNER JOIN specializations_professions sp on spec.id = sp.specialization_id
         INNER JOIN university u on u.id = prog.university_id
WHERE sp.professions_id = $1;`
	rows, err := db.Conn.Query(query, professionID, userID)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var program models.StudyProgramsShort
		err = program.ScanRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, program)
	}
	return list, nil
}

func (db Database) GetAllProgramsForFavouriteProfessions(userID models.UserID) (models.StudyProgramsShortList, error) {
	list := make(models.StudyProgramsShortList, 0)
	query := `SELECT prog.id,
       spec.id,
       spec.name,
       university_id,
       u.name,
       'https://avatars.mds.yandex.net/get-altay/1335362/2a000001649009b50727c93104e9cddcb0cf/XXL',
       exams,
       exists(SELECT fsp.id FROM favourite_study_programs fsp WHERE prog.id = fsp.study_programs_id AND fsp.user_id = $1) as "is_favourite"
FROM study_programs prog
         INNER JOIN specialization spec on spec.id = prog.specialization_id
         INNER JOIN specializations_professions sp on spec.id = sp.specialization_id
         INNER JOIN university u on u.id = prog.university_id
         INNER JOIN favourite_professions fp on sp.professions_id = fp.profession_id
WHERE fp.user_id = $1
GROUP BY prog.id, spec.id, university_id, u.name, exams
ORDER BY count(*) DESC;`
	rows, err := db.Conn.Query(query, userID)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var profession models.StudyProgramsShort
		err := profession.ScanRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, profession)
	}
	return list, nil
}

func (db Database) GetStudyProgramByID(programID int, userId models.UserID) (models.StudyProgramsDetails, error) {
	program := models.StudyProgramsDetails{}
	query := `SELECT prog.id,
       spec.id,
       spec.name,
       spec.description,
       spec.disciplines,
       university_id,
       u.name,
       u.address,
       u.image,
       exams,
       prog.score_budget,
       prog.score_contract,
       prog.contract_amount,
       exists(SELECT fsp.id FROM favourite_study_programs fsp WHERE prog.id = fsp.study_programs_id AND fsp.user_id = $2) as "is_favourite"
FROM study_programs prog
         INNER JOIN specialization spec on spec.id = prog.specialization_id
         INNER JOIN university u on u.id = prog.university_id
WHERE prog.id = $1;`
	row := db.Conn.QueryRow(query, programID, userId)
	switch err := program.ScanRow(row); err {
	case sql.ErrNoRows:
		return program, ErrNoMatch
	default:
		return program, err
	}
}

func (db Database) AddStudyProgramFavourite(userID models.UserID, programID int) error {
	rows := db.Conn.QueryRow("INSERT INTO favourite_study_programs(study_programs_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;",
		programID, userID)
	return rows.Err()
}

func (db Database) DeleteStudyProgramFavourite(userID models.UserID, programID int) error {
	rows := db.Conn.QueryRow("DELETE FROM favourite_study_programs WHERE user_id=$1 AND study_programs_id=$2;",
		userID, programID)
	return rows.Err()
}
