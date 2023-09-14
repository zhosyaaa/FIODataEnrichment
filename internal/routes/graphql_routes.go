package routes

//
//func SetupGraphQLRoutes(router *http.ServeMux, gc *controllers.GraphQLController) {
//	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
//		Query:    gc.RootQuery,    // Укажите корневой запрос
//		Mutation: gc.RootMutation, // Укажите корневую мутацию (если есть)
//	})
//
//	h := handler.New(&handler.Config{
//		Schema:   &schema,
//		Pretty:   true, // Опционально: сделает вывод более читаемым
//		GraphiQL: true, // Опционально: включает интерфейс GraphiQL для отладки
//	})
//
//	router.Handle("/graphql", h)
//}
