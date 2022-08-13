package models

import (
	"net/http"
)

type TestQuestion struct {
	ID     int    `json:"id"`
	TestID int    `json:"test_id"`
	Text   string `json:"text"`
}

func (q *TestQuestion) ScanRow(row ScannedRow) error {
	return row.Scan(
		&q.ID,
		&q.TestID,
		&q.Text,
	)
}

type TestAnswerUserData struct {
	QuestionID int `json:"question_id"`
	Answer     int `json:"answer"`
}

func (u *TestAnswerUserData) Bind(*http.Request) error {
	return nil
}

type TestAnswerUserDataList []TestAnswerUserData

type TestAnswerData struct {
	TestAnswerUserData
	UserID int `json:"user_id"`
}

type TestAnswer struct {
	ID int `json:"id"`
	TestAnswerData
}

func (a *TestAnswer) ScanRow(row ScannedRow) error {
	return row.Scan(
		&a.ID,
		&a.UserID,
		&a.QuestionID,
		&a.Answer,
	)
}

type TestAnswerList []TestAnswer

func (*TestAnswerList) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type Test struct {
	Questions []TestQuestion `json:"questions"`
}

func (*Test) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type TestAnswerCount int

func (*TestAnswerCount) Render(http.ResponseWriter, *http.Request) error {
	return nil
}
