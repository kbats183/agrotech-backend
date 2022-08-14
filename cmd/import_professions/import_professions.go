package main

import (
	"encoding/json"
	"fmt"
	"githab.com/kbats183/argotech/backend/db"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
)

const apiPath = "https://api.hh.ru"

func apiRequestGet(method string, out any) error {
	requestURL := fmt.Sprintf("%s/%s", apiPath, method)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("client: could not create request: %s\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("client: error making http request: %s\n", err)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("client: could not read response body: %s\n", err)
	}
	err = json.Unmarshal(resBody, out)
	if err != nil {
		return fmt.Errorf("failed to unmarsual areas: %s\n", err)
	}
	return nil
}

func loadRegions() (Regions, []Area, error) {
	var areas []Area
	err := apiRequestGet("areas", &areas)
	if err != nil {
		return nil, nil, err
	}

	var regions Regions
	var regionAreas []Area
	for _, area := range areas[0].Areas {
		if area.Name == "Киреевка" {
			continue
		}
		regionID, err := strconv.Atoi(area.ID)
		if err != nil {
			return nil, nil, err
		}
		regions = append(regions, Region{ID: regionID, Name: area.Name})
		regionAreas = append(regionAreas, area)
	}
	sort.Sort(regions)
	return regions, areas, nil
}

func loadVacancies(query string, areaID int) []Vacancies {
	requestMethod := fmt.Sprintf("vacancies?per_page=100&text=%s", url.QueryEscape(query))
	if areaID != 0 {
		requestMethod = fmt.Sprintf("%s&area=%d", requestMethod, areaID)
	}

	var vacancies VacanciesResponse
	err := apiRequestGet(requestMethod, &vacancies)
	if err != nil {
		panic(err)
	}
	return vacancies.Items
}

func findRegionDfs(id string, region string, areas []Area) *string {
	for _, area := range areas {
		if area.ID == id {
			return &region
		}
		if q := findRegionDfs(id, area.ID, area.Areas); q != nil {
			return q
		}

	}
	return nil
}

func findRegion(id string, areas []Area) *string {
	for _, area := range areas {
		if area.ID == id {
			return &id
		}
		if q := findRegionDfs(id, area.ID, area.Areas); q != nil {
			return q
		}
	}
	return nil
}

func saveRegionsToFile() {
	regions, _, err := loadRegions()
	if err != nil {
		panic(err)
	}

	f, err := os.Create("regions.json")
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	err = json.NewEncoder(f).Encode(regions)
	if err != nil {
		panic(err)
	}
}

func putVacancyIntoDB() {
	_, areas, _ := loadRegions()
	profession := "Агроном"
	dbHost, dbPort, dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")
	database, err := db.New(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}
	defer func() { _ = database.Conn.Close() }()

	q := `INSERT INTO vacancies(hh_id, name, url, employer, employer_url, employer_image, responsibility, area)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING ;`
	vacancies := loadVacancies(profession, 0)
	for _, vacancy := range vacancies {
		var region = findRegion(vacancy.Area.Id, areas)
		if region == nil {
			fmt.Println("no area", vacancy.Area)
			region = &vacancy.Area.Id
		}
		_, err := database.Conn.Exec(q,
			vacancy.ID,
			vacancy.Name,
			vacancy.AlternateUrl,
			vacancy.Employer.Name,
			vacancy.Employer.AlternateUrl,
			vacancy.Employer.LogoUrls.Original,
			vacancy.Snippet.Responsibility,
			region,
		)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	//_, areas, err := loadRegions()
	//if err != nil {
	//	panic(err)
	//}
	//reg := findRegion("113", areas)
	//fmt.Println(*reg)
	saveRegionsToFile()
	//putVacancyIntoDB()
}
