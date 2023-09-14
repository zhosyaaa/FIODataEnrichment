package repository

import "TestCase/internal/models"

type PersonRepository interface {
	GetAllPersons() ([]models.Person, error)
	GetPersonByID(id uint) (*models.Person, error)
	CreatePerson(person *models.Person) error
	UpdatePerson(person *models.Person) error
	DeletePerson(id int) error
	FilterPersons(gender string, age int, page int, perPage int) ([]models.Person, error)
}
