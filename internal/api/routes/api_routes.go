package routes

import (
	apiController "TestCase/internal/api/controllers"
	"TestCase/internal/graphql"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
)

type Routes struct {
	apiController.APIController
	GraphQLResolver *graphql.Resolver
}

func NewRoutes(APIController apiController.APIController, graphQLResolver *graphql.Resolver) *Routes {
	return &Routes{APIController: APIController, GraphQLResolver: graphQLResolver}
}

func (r *Routes) SetupAPIRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/persons", r.GetPersons)
		api.GET("/persons/:id", r.GetPerson)
		api.POST("/persons", r.CreatePerson)
		api.PUT("/persons/:id", r.UpdatePerson)
		api.DELETE("/persons/:id", r.DeletePerson)
		api.GET("/persons/filter", r.FilterPersons)
	}
	graphqlHandler := handler.New(&handler.Config{
		Schema:   &graphql.Schema,
		Pretty:   true,
		GraphiQL: true,
	})
	router.POST("/graphql", gin.WrapH(graphqlHandler))

}
