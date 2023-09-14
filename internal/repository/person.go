package repository

import (
	"TestCase/internal/models"
	"gorm.io/gorm"
)

type PersonService struct {
	DB *gorm.DB
}

func NewDatabaseService(db *gorm.DB) PersonRepository {
	return &PersonService{
		DB: db,
	}
}

func (p PersonService) CreatePerson(person *models.Person) error {
	result := p.DB.Create(person)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p PersonService) GetPersonByID(id uint) (*models.Person, error) {
	var person models.Person
	result := p.DB.First(&person, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &person, nil
}

func (p PersonService) UpdatePerson(person *models.Person) error {
	result := p.DB.Save(person)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p PersonService) DeletePersonByID(id uint) error {
	result := p.DB.Delete(&models.Person{}, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return result.Error
	}
	return nil
}

func (p PersonService) FindPeopleByAge(age int) ([]*models.Person, error) {
	var people []*models.Person
	result := p.DB.Where("age = ?", age).Find(&people)
	if result.Error != nil {
		return nil, result.Error
	}
	return people, nil
}

func (p PersonService) GetAllPeople() ([]*models.Person, error) {
	var people []*models.Person
	result := p.DB.Find(&people)
	if result.Error != nil {
		return nil, result.Error
	}
	return people, nil
}

func (p PersonService) CustomQuery() ([]*models.Person, error) {
	var people []*models.Person
	result := p.DB.Where("age > ? AND gender = ?", 30, "Male").Find(&people)
	if result.Error != nil {
		return nil, result.Error
	}
	return people, nil
}
