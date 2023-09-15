package graphql

import (
	"github.com/graphql-go/graphql"
)

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	},
})

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"updateHello": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"newHello": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Логика для обновления "hello"
				return "Updated", nil
			},
		},
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    RootQuery,
	Mutation: RootMutation,
})

var PersonType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Person",
	Fields: graphql.Fields{
		"Name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the person",
		},
		"Surname": &graphql.Field{
			Type:        graphql.String,
			Description: "Surname of the person",
		},
		"Patronymic": &graphql.Field{
			Type:        graphql.String,
			Description: "Patronymic of the person",
		},
		"Age": &graphql.Field{
			Type:        graphql.Int,
			Description: "Age of the person",
		},
		"Gender": &graphql.Field{
			Type:        graphql.String,
			Description: "Gender of the person",
		},
		"Nationality": &graphql.Field{
			Type:        graphql.String,
			Description: "Nationality of the person",
		},
	},
})
