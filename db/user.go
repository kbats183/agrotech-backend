package db

import (
	"database/sql"
	"githab.com/kbats183/argotech/backend/models"
)

func (db Database) GetAllUsers() (models.UserList, error) {
	list := make(models.UserList, 0)
	rows, err := db.Conn.Query("SELECT * FROM users ORDER BY id")
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
	query := `SELECT * FROM users WHERE "login" = $1;`
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
	query := `UPDATE users SET last_name=$2, first_name=$3, patronymic=$4, step=$5, school_class=$6, school_name=$7 WHERE login=$1 RETURNING *;`
	row := db.Conn.QueryRow(
		query,
		auth,
		userData.LastName,
		userData.FirstName,
		userData.Patronymic,
		userData.Step,
		userData.SchoolClass,
		userData.SchoolName)
	err := user.ScanRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNoMatch
		}
	}
	return user, err
}
