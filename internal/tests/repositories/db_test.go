package repositories

import (
	"TestCase/internal/models"
	"TestCase/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func setupTestDB() (*gorm.DB, error) {
	dsn := "postgresql://postgres:1079@localhost:5432/TestCase?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func clearTestDB(db *gorm.DB) {
	db.Exec("DELETE from people")
}

func TestPersonRepository_GetAllPersons(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Error setting up test database: %v", err)
	}
	clearTestDB(db)

	repo := repository.NewPersonRepository(db)

	testPersons := []models.Person{
		{Name: "John", Surname: "Doe", Gender: "Male", Age: 30},
		{Name: "Jane", Surname: "Doe", Gender: "Female", Age: 25},
	}
	for _, person := range testPersons {
		if err := repo.CreatePerson(&person); err != nil {
			t.Fatalf("Error creating test person: %v", err)
		}
	}

	persons, err := repo.GetAllPersons()
	if err != nil {
		t.Fatalf("Error getting all persons: %v", err)
	}

	assert.Equal(t, len(testPersons), len(persons))
}

func TestPersonRepository_GetPersonByID(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Error setting up test database: %v", err)
	}
	clearTestDB(db)

	repo := repository.NewPersonRepository(db)

	testPerson := &models.Person{
		Name:    "John",
		Surname: "Doe",
		Gender:  "Male",
		Age:     30,
	}

	err = repo.CreatePerson(testPerson)
	if err != nil {
		t.Fatalf("Error creating test person: %v", err)
	}

	foundPerson, err := repo.GetPersonByID(testPerson.ID)
	if err != nil {
		t.Fatalf("Error getting person by ID: %v", err)
	}

	assert.NotNil(t, foundPerson)
	assert.Equal(t, testPerson.Name, foundPerson.Name)
	assert.Equal(t, testPerson.Surname, foundPerson.Surname)
}

func TestPersonRepository_CreatePerson(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Error setting up test database: %v", err)
	}

	repo := repository.NewPersonRepository(db)

	testPerson := &models.Person{
		Name:    "John",
		Surname: "Doe",
		Gender:  "Male",
		Age:     30,
	}

	err = repo.CreatePerson(testPerson)
	if err != nil {
		t.Fatalf("Error creating person: %v", err)
	}

	foundPerson, err := repo.GetPersonByID(testPerson.ID)
	if err != nil {
		t.Fatalf("Error getting person by ID: %v", err)
	}

	assert.NotNil(t, foundPerson)
	assert.Equal(t, testPerson.Name, foundPerson.Name)
	assert.Equal(t, testPerson.Surname, foundPerson.Surname)
	assert.Equal(t, testPerson.Gender, foundPerson.Gender)
	assert.Equal(t, testPerson.Age, foundPerson.Age)
}

func TestPersonRepository_UpdatePerson(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Error setting up test database: %v", err)
	}
	clearTestDB(db)

	repo := repository.NewPersonRepository(db)

	testPerson := &models.Person{
		Name:    "John",
		Surname: "Doe",
		Gender:  "Male",
		Age:     30,
	}

	err = repo.CreatePerson(testPerson)
	if err != nil {
		t.Fatalf("Error creating test person: %v", err)
	}

	testPerson.Name = "UpdatedJohn"
	testPerson.Surname = "UpdatedDoe"

	err = repo.UpdatePerson(testPerson)
	if err != nil {
		t.Fatalf("Error updating person: %v", err)
	}

	foundPerson, err := repo.GetPersonByID(testPerson.ID)
	if err != nil {
		t.Fatalf("Error getting person by ID: %v", err)
	}

	assert.NotNil(t, foundPerson)
	assert.Equal(t, testPerson.Name, foundPerson.Name)
	assert.Equal(t, testPerson.Surname, foundPerson.Surname)
}

func TestPersonRepository_DeletePerson(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Error setting up test database: %v", err)
	}
	clearTestDB(db)
	repo := repository.NewPersonRepository(db)

	testPerson := &models.Person{
		Name:    "John",
		Surname: "Doe",
		Gender:  "Male",
		Age:     30,
	}

	err = repo.CreatePerson(testPerson)
	if err != nil {
		t.Fatalf("Error creating person: %v", err)
	}

	err = repo.DeletePerson(int(testPerson.ID))
	if err != nil {
		t.Fatalf("Error deleting person: %v", err)
	}

	foundPerson, _ := repo.GetPersonByID(testPerson.ID)

	assert.Nil(t, foundPerson)
}

func TestPersonRepository_FilterPersons(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Error setting up test database: %v", err)
	}
	clearTestDB(db)
	repo := repository.NewPersonRepository(db)

	testPersons := []models.Person{
		{Name: "John", Surname: "Doe", Gender: "Male", Age: 30},
		{Name: "Jane", Surname: "Doe", Gender: "Female", Age: 25},
		{Name: "Alice", Surname: "Smith", Gender: "Female", Age: 35},
	}

	for _, person := range testPersons {
		err := repo.CreatePerson(&person)
		if err != nil {
			t.Fatalf("Error creating test person: %v", err)
		}
	}

	filteredPersons, err := repo.FilterPersons("Male", 30, 1, 10)
	if err != nil {
		t.Fatalf("Error filtering persons: %v", err)
	}
	assert.Equal(t, 1, len(filteredPersons))
}
