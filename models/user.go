package models

import (
	"net/http"
)

type UserID int

func (*UserID) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type UserAuth string

func (u *UserAuth) Bind(*http.Request) error {
	return nil
}

type UserProfile struct {
	LastName               string  `json:"last_name"`
	FirstName              string  `json:"first_name"`
	Patronymic             string  `json:"patronymic"`
	Region                 *int    `json:"region"`
	Step                   *int    `json:"step"`
	SchoolClass            *int    `json:"school_class"`
	SchoolName             *string `json:"school_name"`
	UniversityName         *string `json:"university_name"`
	UniversityStudyProgram *string `json:"university_study_program"`
	UniversityProfession   *int    `json:"university_profession"`
}

func (*UserProfile) Bind(*http.Request) error {
	return nil
}

type UserData struct {
	Login string `json:"login"`
	UserProfile
}

func (u *UserData) Bind(*http.Request) error {
	return nil
}

type User struct {
	UserData
	ID        UserID `json:"id"`
	CreatedAt string `json:"created_at"`
}

func (user *User) ScanRow(row ScannedRow) error {
	return row.Scan(
		&user.ID,
		&user.Login,
		&user.LastName,
		&user.FirstName,
		&user.Patronymic,
		&user.Region,
		&user.Step,
		&user.SchoolClass,
		&user.SchoolName,
		&user.UniversityName,
		&user.UniversityStudyProgram,
		&user.UniversityProfession,
		&user.CreatedAt)
}

func (*User) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type UserList []User

func (UserList) Render(http.ResponseWriter, *http.Request) error {
	return nil
}
