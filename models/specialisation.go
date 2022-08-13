package models

import (
	"encoding/json"
	"net/http"
)

type ExamsList []uint

func (*ExamsList) Bind(*http.Request) error {
	return nil
}

type StudyProgramsShort struct {
	ID               int       `json:"id"`
	SpecialisationID int       `json:"specialisation_id"`
	UniversityID     int       `json:"university_id"`
	UniversityName   string    `json:"university_name"`
	Exams            ExamsList `json:"exams"`
}

func (s *StudyProgramsShort) ScanRow(row ScannedRow) error {
	var exam []byte
	err := row.Scan(&s.ID, &s.SpecialisationID, &s.UniversityID, &s.UniversityName, &exam)
	if err != nil {
		return err
	}
	return json.Unmarshal(exam, &s.Exams)
}

type StudyProgramsShortList []StudyProgramsShort

func (*StudyProgramsShortList) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type StudyProgramsUniversityAndSpecialisation struct {
	StudyProgramsShort
	UniversityImage    string `json:"university_image"`
	SpecialisationName string `json:"specialisation_name"`
	IsFavourite        bool   `json:"is_favourite"`
}

func (s *StudyProgramsUniversityAndSpecialisation) ScanRow(row ScannedRow) error {
	var exam []byte
	err := row.Scan(
		&s.ID,
		&s.SpecialisationID,
		&s.SpecialisationName,
		&s.UniversityID,
		&s.UniversityName,
		&s.UniversityImage,
		&exam,
		&s.IsFavourite)
	if err != nil {
		return err
	}
	return json.Unmarshal(exam, &s.Exams)
}

func (*StudyProgramsUniversityAndSpecialisation) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type StudyProgramsUniversityAndSpecialisationList []StudyProgramsUniversityAndSpecialisation

func (*StudyProgramsUniversityAndSpecialisationList) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type StudyProgramsDetails struct {
	StudyProgramsUniversityAndSpecialisation
	UniversityAddress         string `json:"university_address"`
	SpecialisationDescription string `json:"specialisation_description"`
	SpecialisationDisciplines string `json:"specialisation_disciplines"`
	ScoreBudget               int    `json:"score_budget"`
	ScoreContract             int    `json:"score_contract"`
	ContractAmount            int    `json:"contract_amount"`
}

func (s *StudyProgramsDetails) ScanRow(row ScannedRow) error {
	var exam []byte
	err := row.Scan(
		&s.ID,
		&s.SpecialisationID,
		&s.SpecialisationName,
		&s.SpecialisationDescription,
		&s.SpecialisationDisciplines,
		&s.UniversityID,
		&s.UniversityName,
		&s.UniversityAddress,
		&s.UniversityImage,
		&exam,
		&s.ScoreBudget,
		&s.ScoreContract,
		&s.ContractAmount,
		&s.IsFavourite)
	if err != nil {
		return err
	}
	return json.Unmarshal(exam, &s.Exams)
}

func (*StudyProgramsDetails) Render(http.ResponseWriter, *http.Request) error {
	return nil
}
