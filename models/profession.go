package models

import (
	"net/http"
)

type ProfessionShortInfo struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ShortDescription *string `json:"short_description"`
}

func (p *ProfessionShortInfo) ScanRow(row ScannedRow) error {
	return row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.ShortDescription,
	)
}

type ProfessionShortInfoList []ProfessionShortInfo

func (*ProfessionShortInfoList) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type ProfessionShortInfoWithRating struct {
	ProfessionShortInfo
	IsFavourite bool `json:"is_favourite"`
}

func (p *ProfessionShortInfoWithRating) ScanRow(row ScannedRow) error {
	return row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.ShortDescription,
		&p.IsFavourite,
	)
}

type ProfessionShortInfoWithRatingList []ProfessionShortInfoWithRating

func (*ProfessionShortInfoWithRatingList) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type Profession struct {
	ProfessionShortInfo
	Tasks          string  `json:"tasks"`
	RequiredSkills string  `json:"required_skills"`
	Relevance      *string `json:"relevance"`
}

func (p *Profession) ScanRow(row ScannedRow) error {
	return row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.ShortDescription,
		&p.Tasks,
		&p.RequiredSkills,
		&p.Relevance)
}

func (*Profession) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type ProfessionList []Profession

func (*ProfessionList) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type ProfessionWithRating struct {
	Profession
	IsFavourite bool `json:"is_favourite"`
}

func (p *ProfessionWithRating) ScanRow(row ScannedRow) error {
	return row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.ShortDescription,
		&p.Tasks,
		&p.RequiredSkills,
		&p.Relevance,
		&p.IsFavourite,
	)
}

type ProfessionWithRatingList []ProfessionWithRating

func (*ProfessionWithRatingList) Render(http.ResponseWriter, *http.Request) error {
	return nil
}
