package models

import (
	"net/http"
)

type Vacancy struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Area           string `json:"area"`
	Employer       string `json:"employer"`
	EmployerUrl    string `json:"employer_url"`
	EmployerImage  string `json:"employer_image"`
	Responsibility string `json:"responsibility"`
	Url            string `json:"url"`
}

func (v *Vacancy) ScanRow(row ScannedRow) error {
	var hhID int
	return row.Scan(
		&v.ID,
		&hhID,
		&v.Name,
		&v.Url,
		&v.Employer,
		&v.EmployerUrl,
		&v.EmployerImage,
		&v.Responsibility,
		&v.Area,
	)
}

type VacancyList []Vacancy

func (*VacancyList) Render(http.ResponseWriter, *http.Request) error {
	return nil
}
