package controllers

import (
	"TestCase/internal/models"
	"TestCase/internal/redis"
	"TestCase/internal/repository"
	"TestCase/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type APIController struct {
	PersonRepository  repository.PersonRepository
	CacheClient       redis.CacheClient
	enrichmentService *services.EnrichmentService
}

func NewAPIController(personRepository repository.PersonRepository, cacheClient redis.CacheClient, enrichmentService *services.EnrichmentService) *APIController {
	return &APIController{PersonRepository: personRepository, CacheClient: cacheClient, enrichmentService: enrichmentService}
}

func (ac *APIController) GetPersons(c *gin.Context) {
	persons, err := ac.PersonRepository.GetAllPersons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, persons)
}

func (ac *APIController) CreatePerson(c *gin.Context) {
	var input *models.Input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var person models.Person
	person.Name = input.Name
	person.Surname = input.Surname
	person.Patronymic = input.Patronymic

	fio := fmt.Sprintf("%s %s %s", input.Name, input.Surname, input.Patronymic)
	ac.enrichmentService.FIOChannel <- fio

	c.JSON(http.StatusCreated, person)
}

func (ac *APIController) GetPerson(c *gin.Context) {
	idParam := c.Param("id")
	personID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
		return
	}

	person, err := ac.CacheClient.GetPerson(uint(personID))
	if err == nil {
		c.JSON(http.StatusOK, person)
		return
	}

	person, err = ac.PersonRepository.GetPersonByID(uint(personID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	if err := ac.CacheClient.SetPerson(*person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

func (ac *APIController) UpdatePerson(c *gin.Context) {
	idParam := c.Param("id")
	personID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
		return
	}

	var updatedPerson models.Person
	if err := c.BindJSON(&updatedPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingPerson, err := ac.PersonRepository.GetPersonByID(uint(personID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	existingPerson.Name = updatedPerson.Name
	existingPerson.Surname = updatedPerson.Surname
	existingPerson.Patronymic = updatedPerson.Patronymic
	existingPerson.Age = updatedPerson.Age
	existingPerson.Gender = updatedPerson.Gender
	existingPerson.Nationality = updatedPerson.Nationality

	if err := ac.PersonRepository.UpdatePerson(existingPerson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ac.CacheClient.SetPerson(*existingPerson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingPerson)
}

func (ac *APIController) DeletePerson(c *gin.Context) {
	idParam := c.Param("id")
	personID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
		return
	}

	if err := ac.PersonRepository.DeletePerson(personID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ac.CacheClient.DeletePerson(uint(personID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (ac *APIController) FilterPersons(c *gin.Context) {
	gender := c.Query("gender")
	ageParam := c.Query("age")
	pageParam := c.DefaultQuery("page", "1")
	perPageParam := c.DefaultQuery("per_page", "10")

	age, err := strconv.Atoi(ageParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age parameter"})
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	perPage, err := strconv.Atoi(perPageParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid per_page parameter"})
		return
	}

	filteredPersons, err := ac.PersonRepository.FilterPersons(gender, age, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, filteredPersons)
}
