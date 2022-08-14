package models

import (
	"net/http"
)

type University struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Image   string `json:"image"`
}

func (*University) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

func (u *University) ScanRow(row ScannedRow) error {
	return row.Scan(&u.ID, &u.Name, &u.Address, &u.Image)
}

type UniversityList []University

func (*UniversityList) Render(http.ResponseWriter, *http.Request) error {
	return nil
}
