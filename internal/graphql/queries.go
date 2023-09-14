package graphql

//
//var RootQuery = graphql.NewObject(graphql.ObjectConfig{
//	Name: "RootQuery",
//	Fields: graphql.Fields{
//		"getPeopleWithFiltersAndPagination": &graphql.Field{
//			Type:       , // Укажите тип данных, который этот запрос возвращает,
//			Description: "Получить людей с фильтрами и пагинацией",
//			Args: graphql.FieldConfigArgument{
//			"filters": &graphql.ArgumentConfig{
//			Type: ,// Укажите тип аргумента фильтров,
//		},
//			"page": &graphql.ArgumentConfig{
//			Type: graphql.Int,
//		},
//			"perPage": &graphql.ArgumentConfig{
//			Type: graphql.Int,
//		},
//		},
//			Resolve: gc.GetPeopleWithFiltersAndPagination, // Укажите метод контроллера для обработки запроса
//		},
//	},
//})
