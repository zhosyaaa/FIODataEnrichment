package repository

import "TestCase/internal/models"

type PersonRepository interface {
	CreatePerson(person *models.Person) error
	GetPersonByID(id uint) (*models.Person, error)
	UpdatePerson(person *models.Person) error
	DeletePersonByID(id uint) error
	FindPeopleByAge(age int) ([]*models.Person, error)
	GetAllPeople() ([]*models.Person, error)
	CustomQuery() ([]*models.Person, error)
}
