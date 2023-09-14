package services

import (
	"TestCase/internal/models"
	"TestCase/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIResponse struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

type EnrichmentService struct {
	AgifyAPIClient       *http.Client
	GenderizeAPIClient   *http.Client
	NationalizeAPIClient *http.Client
	DatabaseService      *repository.PersonService
}

func NewEnrichmentService(agifyClient, genderizeClient, nationalizeClient *http.Client, dbService *repository.PersonService) *EnrichmentService {
	return &EnrichmentService{
		AgifyAPIClient:       agifyClient,
		GenderizeAPIClient:   genderizeClient,
		NationalizeAPIClient: nationalizeClient,
		DatabaseService:      dbService,
	}
}

func (es *EnrichmentService) EnrichAndSaveFIO(fioStream <-chan string) {
	for fio := range fioStream {
		// Обработка ФИО из потока
		person := &models.Person{
			Name: fio,
		}

		es.enrichPersonData(person)

		// Сохранение данных в базе данных
		if err := es.DatabaseService.CreatePerson(person); err != nil {
			fmt.Printf("Ошибка при сохранении данных в базу данных: %v\n", err)
		}
	}
}

func (es *EnrichmentService) enrichPersonData(person *models.Person) {
	// Обогащение данных о возрасте
	agifyURL := fmt.Sprintf("https://api.agify.io/?name=%s", person.Name)
	agifyResponse, err := es.fetchAPI(agifyURL)
	if err != nil {
		fmt.Printf("Ошибка при запросе возраста: %v\n", err)
		return
	}
	person.Age = agifyResponse.Age

	genderizeURL := fmt.Sprintf("https://api.genderize.io/?name=%s", person.Name)
	genderizeResponse, err := es.fetchAPI(genderizeURL)
	if err != nil {
		fmt.Printf("Ошибка при запросе пола: %v\n", err)
		return
	}
	person.Gender = genderizeResponse.Gender

	// Обогащение данных о национальности
	nationalizeURL := fmt.Sprintf("https://api.nationalize.io/?name=%s", person.Name)
	nationalizeResponse, err := es.fetchAPI(nationalizeURL)
	if err != nil {
		fmt.Printf("Ошибка при запросе национальности: %v\n", err)
		return
	}
	if len(nationalizeResponse.Country) > 0 {
		person.Nationality = nationalizeResponse.Country[0].CountryID
	}
}

func (es *EnrichmentService) fetchAPI(apiURL string) (*APIResponse, error) {
	resp, err := es.AgifyAPIClient.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
