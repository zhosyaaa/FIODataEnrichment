package controllers

import (
	"TestCase/internal/models"
	"TestCase/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type APIController struct {
	PersonRepository repository.PersonRepositoryImpl
}

func NewAPIController(personRepo repository.PersonRepositoryImpl) *APIController {
	return &APIController{
		PersonRepository: personRepo,
	}
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
	var person models.Person
	if err := c.BindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.PersonRepository.CreatePerson(&person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, person)
}

func (ac *APIController) GetPerson(c *gin.Context) {
	idParam := c.Param("id")
	personID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
		return
	}

	person, err := ac.PersonRepository.GetPersonByID(uint(personID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
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

	// Выполните обновление данных в существующем объекте `existingPerson`
	existingPerson.Name = updatedPerson.Name
	existingPerson.Surname = updatedPerson.Surname
	// Обновите остальные поля по аналогии

	if err := ac.PersonRepository.UpdatePerson(existingPerson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingPerson)
}

func (ac *APIController) DeletePerson(c *gin.Context) {
	// Удаление физического лица по идентификатору
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

	c.JSON(http.StatusNoContent, nil)
}

func (ac *APIController) FilterPersons(c *gin.Context) {
	// Пример: /api/persons/filter?gender=male&age=30&page=1&per_page=10
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

	// Здесь вы можете использовать ваш репозиторий для фильтрации и пагинации данных
	filteredPersons, err := ac.PersonRepository.FilterPersons(gender, age, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, filteredPersons)
}
