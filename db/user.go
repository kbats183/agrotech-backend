package db

import (
	"database/sql"
	"githab.com/kbats183/argotech/backend/models"
)

func (db Database) GetAllUsers() (models.UserList, error) {
	list := make(models.UserList, 0)
	query := `SELECT id, 
       login,
       last_name,
       first_name,
       patronymic,
       region,
       step,
       school_class,
       school_name,
       university_name,
       university_study_program,
       university_profession,
       created_at
FROM users ORDER BY id`
	rows, err := db.Conn.Query(query)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var user models.User
		err := user.ScanRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, user)
	}
	return list, nil
}
func (db Database) AddUser(user *models.User) error {
	var id models.UserID
	var createdAt string
	query := `INSERT INTO users (login, last_name, first_name, patronymic) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, user.Login, user.LastName, user.FirstName, user.Patronymic).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	user.ID = id
	user.CreatedAt = createdAt
	return nil
}
func (db Database) GetUserByAuth(auth string) (models.User, error) {
	user := models.User{}
	query := `SELECT id, 
       login,
       last_name,
       first_name,
       patronymic,
       region,
       step,
       school_class,
       school_name,
       university_name,
       university_study_program,
       university_profession,
       created_at 
FROM users WHERE "login" = $1;`
	row := db.Conn.QueryRow(query, auth)
	switch err := user.ScanRow(row); err {
	case sql.ErrNoRows:
		return user, ErrNoMatch
	default:
		return user, err
	}
}
func (db Database) DeleteUser(userId int) error {
	query := `DELETE FROM users WHERE id = $1;`
	_, err := db.Conn.Exec(query, userId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}
func (db Database) UpdateUser(userID int, userData models.UserData) (models.User, error) {
	user := models.User{UserData: userData}
	query := `UPDATE users SET login=$2, last_name=$3, first_name=$4, patronymic=$5 WHERE id=$1 RETURNING id, created_at;`
	err := db.Conn.QueryRow(query, userID, userData.Login, userData.LastName, userData.FirstName, user.Patronymic).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNoMatch
		}
		return user, err
	}
	return user, nil
}
func (db Database) UpdateUserProfile(auth string, userData models.UserProfile) (models.User, error) {
	var user models.User
	query := `UPDATE users SET last_name=$2, first_name=$3, patronymic=$4, region=$5, step=$6, school_class=$7, school_name=$8, university_name=$9, university_study_program=$10, university_profession=$11 WHERE login=$1 RETURNING 
       id,
       login,
       last_name,
       first_name,
       patronymic,
       region,
       step,
       school_class,
       school_name,
       university_name,
       university_study_program,
       university_profession,
       created_at ;`
	row := db.Conn.QueryRow(
		query,
		auth,
		userData.LastName,
		userData.FirstName,
		userData.Patronymic,
		userData.Region,
		userData.Step,
		userData.SchoolClass,
		userData.SchoolName,
		userData.UniversityName,
		userData.UniversityStudyProgram,
		userData.UniversityProfession,
	)
	err := user.ScanRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNoMatch
		}
	}
	return user, err
}

func (db Database) GetCVByLogin(login string) (models.UserCV, error) {
	user := models.UserCV{}
	query := `SELECT
       u.login,
       u.last_name,
       u.first_name,
       u.patronymic,
       u.date_of_birth,
       u.address,
       u.contact_data,
       u.school_name,
       u.school_begin_year,
       u.school_end_year,
       u.university_name,
       u.university_study_program,
       u.university_begin_year,
       u.university_end_year,
       u.skills
FROM users u
WHERE u.login = $1;`
	row := db.Conn.QueryRow(query, login)
	switch err := user.ScanRow(row); err {
	case sql.ErrNoRows:
		return user, ErrNoMatch
	default:
		return user, err
	}
}

func (db Database) UpdateCV(login string, d models.UserCV) error {
	query := `UPDATE users
SET last_name                = $2,
    first_name               = $3,
    patronymic               = $4,
    date_of_birth            = $5,
    address                  = $6,
    contact_data             = $7,
    school_name              = $8,
    school_begin_year        = $9,
    school_end_year          = $10,
    university_name          = $11,
    university_study_program = $12,
    university_begin_year    = $13,
    university_end_year      = $14,
    skills                   = $15
WHERE login = $1;`
	_, err := db.Conn.Exec(query,
		login,
		d.LastName,
		d.FirstName,
		d.Patronymic,
		d.DateOfBirth,
		d.Address,
		d.ContactData,
		d.SchoolName,
		d.SchoolBeginYear,
		d.SchoolEndYear,
		d.UniversityName,
		d.UniversityStudyProgram,
		d.UniversityBeginYear,
		d.UniversityEndYear,
		d.Skills,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoMatch
		}
		return err
	}
	return nil
}
