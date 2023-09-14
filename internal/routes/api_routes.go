package routes

import (
	apiController "TestCase/internal/controllers"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	apiController.APIController
}

func NewRoutes(APIController apiController.APIController) *Routes {
	return &Routes{APIController: APIController}
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
}
