package models

import (
	"net/http"
)

type Profession struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ShortDescription *string `json:"short_description"`
	Tasks            string  `json:"tasks"`
	RequiredSkills   string  `json:"required_skills"`
	Relevance        *string `json:"relevance"`
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
