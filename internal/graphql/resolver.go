package graphql

import (
	"TestCase/internal/api/services"
	"TestCase/internal/db"
	"TestCase/internal/models"
	"TestCase/internal/repository"
	"context"
	"fmt"
	"net/http"
)

type Resolver struct {
	PersonRepository  repository.PersonRepository
	EnrichmentService services.EnrichmentService
}

func NewResolver(personRepository repository.PersonRepository) *Resolver {
	return &Resolver{PersonRepository: personRepository}
}

func (r *Resolver) resolvePersons(ctx context.Context) ([]models.Person, error) {
	persons, err := r.PersonRepository.GetAllPersons()
	if err != nil {
		return nil, err
	}
	return persons, nil
}

func (r *Resolver) resolvePerson(ctx context.Context, ID int) (*models.Person, error) {
	person, err := r.PersonRepository.GetPersonByID(uint(ID))
	if err != nil {
		return nil, err
	}
	return person, nil
}

func (r *Resolver) resolveFilteredPersons(ctx context.Context, args struct {
	Gender  *string
	Age     *int
	Page    *int
	PerPage *int
}) ([]models.Person, error) {
	filteredPersons, err := r.PersonRepository.FilterPersons(
		*args.Gender, *args.Age, *args.Page, *args.PerPage)
	if err != nil {
		return nil, err
	}
	return filteredPersons, nil
}

func (r *Resolver) resolveCreatePerson(ctx context.Context, Input models.Input) (*models.Person, error) {
	personService := &repository.PersonRepositoryImpl{DB: db.DB}
	r.EnrichmentService = *services.NewEnrichmentService(
		&http.Client{},
		&http.Client{},
		&http.Client{},
		personService,
		make(chan string),
		nil,
	)

	newPerson := &models.Person{
		Name:       Input.Name,
		Surname:    Input.Surname,
		Patronymic: Input.Patronymic,
	}
	fio := fmt.Sprintf("%s %s %s", newPerson.Name, newPerson.Surname, newPerson.Patronymic)
	go r.EnrichmentService.EnrichAndSaveFIO()
	r.EnrichmentService.FIOChannel <- fio
	return newPerson, nil
}

func (r *Resolver) resolveUpdatePerson(ctx context.Context, ID int, Input models.Person) (*models.Person, error) {
	person, err := r.PersonRepository.GetPersonByID(uint(ID))
	if err != nil {
		return nil, err
	}

	person.Name = Input.Name
	person.Surname = Input.Surname
	person.Patronymic = Input.Patronymic
	person.Age = Input.Age
	person.Gender = Input.Gender
	person.Nationality = Input.Nationality

	if err := r.PersonRepository.UpdatePerson(person); err != nil {
		return nil, err
	}
	return person, nil
}

func (r *Resolver) resolveDeletePerson(ctx context.Context, ID int) (bool, error) {
	if err := r.PersonRepository.DeletePerson(ID); err != nil {
		return false, err
	}
	return true, nil
}
