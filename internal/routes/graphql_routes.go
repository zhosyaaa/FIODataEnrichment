package routes

//func (r *GraphQLRoutes) SetupGraphQLRoutes(router *gin.Engine) {
//	graphqlRouter := router.Group("/graphql")
//	{
//		graphqlRouter.POST("", r.graphqlController.HandleGraphQL)
//		graphqlRouter.GET("/persons", r.graphqlController.GetPersons)           // Аналогичный маршрут для получения персон
//		graphqlRouter.GET("/persons/:id", r.graphqlController.GetPersonByID)    // Аналогичный маршрут для получения персоны по ID
//		graphqlRouter.POST("/persons", r.graphqlController.CreatePerson)        // Аналогичный маршрут для создания персоны
//		graphqlRouter.PUT("/persons/:id", r.graphqlController.UpdatePerson)     // Аналогичный маршрут для обновления персоны
//		graphqlRouter.DELETE("/persons/:id", r.graphqlController.DeletePerson)  // Аналогичный маршрут для удаления персоны
//		graphqlRouter.GET("/persons/filter", r.graphqlController.FilterPersons) // Аналогичный маршрут для фильтрации персон
//	}
//}
