package repository

import (
	models2 "TestCase/internal/models"
	"gorm.io/gorm"
)

type PersonRepositoryImpl struct {
	DB *gorm.DB
}

func NewPersonRepository(db *gorm.DB) *PersonRepositoryImpl {
	return &PersonRepositoryImpl{
		DB: db,
	}
}

func (r *PersonRepositoryImpl) GetAllPersons() ([]models2.Person, error) {
	var persons []models2.Person
	if err := r.DB.Find(&persons).Error; err != nil {
		return nil, err
	}
	return persons, nil
}

func (r *PersonRepositoryImpl) GetPersonByID(id uint) (*models2.Person, error) {
	var person models2.Person
	if err := r.DB.First(&person, id).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *PersonRepositoryImpl) CreatePerson(person *models2.Person) error {
	if err := r.DB.Create(person).Error; err != nil {
		return err
	}
	return nil
}

func (r *PersonRepositoryImpl) UpdatePerson(person *models2.Person) error {
	if err := r.DB.Save(person).Error; err != nil {
		return err
	}
	return nil
}

func (r *PersonRepositoryImpl) DeletePerson(id int) error {
	if err := r.DB.Delete(&models2.Person{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *PersonRepositoryImpl) FilterPersons(gender string, age int, page int, perPage int) ([]models2.Person, error) {
	var persons []models2.Person
	query := r.DB
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if age > 0 {
		query = query.Where("age = ?", age)
	}
	if err := query.Offset((page - 1) * perPage).Limit(perPage).Find(&persons).Error; err != nil {
		return nil, err
	}
	return persons, nil
}
