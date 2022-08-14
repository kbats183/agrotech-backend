package main

import "strings"

type Area struct {
	ID       string  `json:"id"`
	ParentID *string `json:"parentID"`
	Name     string  `json:"name"`
	Areas    []Area  `json:"areas"`
}

type Region struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Regions []Region

func (r Regions) Len() int {
	return len(r)
}

func (r Regions) Less(i, j int) bool {
	return strings.Compare(r[i].Name, r[j].Name) < 1
}

func (r Regions) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type Vacancies struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AlternateUrl string `json:"alternate_url"`
	Area         struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"area"`
	Employer struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Url      string `json:"url"`
		LogoUrls struct {
			Original string `json:"original"`
		} `json:"logo_urls"`
		AlternateUrl string `json:"alternate_url"`
	} `json:"employer"`
	Snippet struct {
		Requirement    string `json:"requirement"`
		Responsibility string `json:"responsibility"`
	} `json:"snippet"`
}

type VacanciesResponse struct {
	Items []Vacancies `json:"items"`
}
