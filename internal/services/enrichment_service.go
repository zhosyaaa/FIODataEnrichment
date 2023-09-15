package services

import (
	"TestCase/internal/configs"
	"TestCase/internal/models"
	"TestCase/internal/repository"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"net/http"
	"strings"
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
	DatabaseService      *repository.PersonRepositoryImpl
	FIOChannel           chan string
	RedisClient          *redis.Client
}

func NewEnrichmentService(agifyClient, genderizeClient, nationalizeClient *http.Client, dbService *repository.PersonRepositoryImpl) *EnrichmentService {
	return &EnrichmentService{
		AgifyAPIClient:       agifyClient,
		GenderizeAPIClient:   genderizeClient,
		NationalizeAPIClient: nationalizeClient,
		DatabaseService:      dbService,
	}
}
func init() {
	configs.InitRedis()
}

func (es *EnrichmentService) EnrichAndSaveFIO(fioStream <-chan string) {
	for fio := range fioStream {
		cachedData, err := configs.GetFromCache(fio)
		if err == nil {
			fmt.Printf("Cached data used for: %s\n", fio)
			var cachedPerson models.Person
			json.Unmarshal([]byte(cachedData), &cachedPerson)
		} else {
			fioInformation := strings.Split(fio, " ")
			person := &models.Person{
				Name:       fioInformation[0],
				Surname:    fioInformation[1],
				Patronymic: fioInformation[2],
			}
			es.enrichPersonData(person)
			jsonBytes, _ := json.Marshal(person)
			configs.SetInCache(fio, string(jsonBytes))

			if err := es.DatabaseService.CreatePerson(person); err != nil {
				fmt.Printf("Error saving data to the database: %v\n", err)
			}
		}
	}
}

func (es *EnrichmentService) enrichPersonData(person *models.Person) {
	fio := fmt.Sprintf("%s %s %s", person.Name, person.Surname, person.Patronymic)
	agifyURL := fmt.Sprintf("https://api.agify.io/?name=%s", fio)
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
