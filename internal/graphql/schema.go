package graphql

import (
	"TestCase/internal/db"
	"TestCase/internal/models"
	"TestCase/internal/repository"
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log" // Импортируйте zerolog
	"strconv"
)

var (
	personType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Person",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.ID},
			"name":        &graphql.Field{Type: graphql.String},
			"surname":     &graphql.Field{Type: graphql.String},
			"patronymic":  &graphql.Field{Type: graphql.String},
			"age":         &graphql.Field{Type: graphql.Int},
			"gender":      &graphql.Field{Type: graphql.String},
			"nationality": &graphql.Field{Type: graphql.String},
		},
	})

	queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"persons": &graphql.Field{
				Type: graphql.NewList(personType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					pr := repository.PersonRepositoryImpl{
						DB: db.DB,
					}
					resolver := NewResolver(&pr)
					return resolver.resolvePersons(p.Context)
				},
			},
			"person": &graphql.Field{
				Type: personType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(string)
					idValue, err := strconv.Atoi(id)
					if err != nil {
						log.Error().Err(err).Msg("Invalid person ID")
						return nil, errors.New("Invalid person ID")
					}
					pr := repository.PersonRepositoryImpl{
						DB: db.DB,
					}
					resolver := NewResolver(&pr)
					return resolver.resolvePerson(p.Context, idValue)
				},
			},
			"filteredPersons": &graphql.Field{
				Type: graphql.NewList(personType),
				Args: graphql.FieldConfigArgument{
					"gender":  &graphql.ArgumentConfig{Type: graphql.String},
					"age":     &graphql.ArgumentConfig{Type: graphql.Int},
					"page":    &graphql.ArgumentConfig{Type: graphql.Int},
					"perPage": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					gender, _ := p.Args["gender"].(string)
					age, _ := p.Args["age"].(int)
					page, _ := p.Args["page"].(int)
					perPage, _ := p.Args["perPage"].(int)
					pr := repository.PersonRepositoryImpl{
						DB: db.DB,
					}
					resolver := NewResolver(&pr)
					args := struct {
						Gender  *string
						Age     *int
						Page    *int
						PerPage *int
					}{
						Gender:  &gender,
						Age:     &age,
						Page:    &page,
						PerPage: &perPage,
					}
					return resolver.resolveFilteredPersons(p.Context, args)
				},
			},
		},
	})

	mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createPerson": &graphql.Field{
				Type: personType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: InputType},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input, ok := p.Args["input"].(map[string]interface{})
					if !ok {
						log.Error().Msg("Invalid input format")
						return nil, errors.New("Invalid input format")
					}
					name, nameOk := input["name"].(string)
					surname, surnameOk := input["surname"].(string)
					patronymic, patronymicOk := input["patronymic"].(string)

					if !nameOk || !surnameOk || !patronymicOk {
						log.Error().Msg("Name, surname, and patronymic are required fields")
						return nil, errors.New("Name, surname, and patronymic are required fields")
					}
					newPerson := models.Input{
						Name:       name,
						Surname:    surname,
						Patronymic: patronymic,
					}
					pr := repository.PersonRepositoryImpl{
						DB: db.DB,
					}
					resolver := NewResolver(&pr)
					return resolver.resolveCreatePerson(p.Context, newPerson)
				},
			},
			"updatePerson": &graphql.Field{
				Type: personType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: UpdateType},
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					pr := repository.PersonRepositoryImpl{
						DB: db.DB,
					}
					resolver := NewResolver(&pr)
					id, idOK := p.Args["id"].(string)
					input, inputOK := p.Args["input"].(map[string]interface{})
					if !idOK || !inputOK {
						log.Error().Msg("Invalid argument types")
						return nil, errors.New("Invalid argument types")
					}
					idValue, err := strconv.Atoi(id)
					person, err := resolver.resolvePerson(p.Context, idValue)
					if err != nil {
						log.Error().Err(err).Msg("Person not found")
						return nil, err
					}
					if err != nil {
						log.Error().Err(err).Msg("Error updating person")
						return nil, err
					}
					if name, ok := input["name"].(string); ok {
						person.Name = name
					}
					if surname, ok := input["surname"].(string); ok {
						person.Surname = surname
					}
					if patronymic, ok := input["patronymic"].(string); ok {
						person.Patronymic = patronymic
					}
					if age, ok := input["age"].(int); ok {
						person.Age = age
					}
					if gender, ok := input["gender"].(string); ok {
						person.Gender = gender
					}
					if nationality, ok := input["nationality"].(string); ok {
						person.Nationality = nationality
					}
					updatedPerson, err := resolver.resolveUpdatePerson(p.Context, idValue, *person)
					if err != nil {
						log.Error().Err(err).Msg("Error updating person")
						return nil, err
					}

					return updatedPerson, nil
				},
			},
			"deletePerson": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, idOK := p.Args["id"].(string)
					if !idOK {
						log.Error().Msg("Invalid argument types")
						return false, errors.New("Invalid argument types")
					}
					idValue, err := strconv.Atoi(id)
					if err != nil {
						log.Error().Err(err).Msg("Invalid person ID")
						return false, err
					}

					pr := repository.PersonRepositoryImpl{
						DB: db.DB,
					}
					resolver := NewResolver(&pr)
					success, err := resolver.resolveDeletePerson(p.Context, idValue)
					if err != nil {
						log.Error().Err(err).Msg("Error deleting person")
						return false, err
					}

					return success, nil
				},
			},
		},
	})
)

var InputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "Input",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"surname": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"patronymic": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

var UpdateType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "Input",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"surname": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"patronymic": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"age": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"gender": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"nationality": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    queryType,
	Mutation: mutationType,
})
