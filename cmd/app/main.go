package main

import (
	"TestCase/internal/api/controllers"
	"TestCase/internal/api/routes"
	services2 "TestCase/internal/api/services"
	"TestCase/internal/configs"
	"TestCase/internal/db"
	"TestCase/internal/graphql"
	"TestCase/internal/redis"
	"TestCase/internal/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
}

func main() {
	db.InitDatabase()
	redisClient := configs.InitRedis()
	if redisClient == nil {
		log.Fatal().Msg("Failed to initialize Redis client")
		return
	}
	defer redisClient.Close()

	router := gin.Default()
	personRepository := repository.NewPersonRepository(db.DB)
	enrichmentService := services2.NewEnrichmentService(
		&http.Client{},
		&http.Client{},
		&http.Client{},
		personRepository,
		make(chan string),
		redisClient,
	)
	go enrichmentService.EnrichAndSaveFIO()

	kafkaReader := configs.InitKafkaReader()
	kafkaWriter := configs.InitKafkaWriter()
	kafkaService := services2.NewKafkaService(kafkaReader, kafkaWriter, enrichmentService)
	go kafkaService.ConsumeMessages()

	GraphQLResolver := graphql.NewResolver(personRepository)
	cacheClient := redis.NewCacheClient(redisClient)
	apiController := controllers.NewAPIController(personRepository, *cacheClient, enrichmentService)
	apiRoutes := routes.NewRoutes(*apiController, GraphQLResolver)
	apiRoutes.SetupAPIRoutes(router)

	serverPort := "8080"
	server := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	fmt.Printf("Server is running on port %s...\n", serverPort)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Server error")
	}
}
