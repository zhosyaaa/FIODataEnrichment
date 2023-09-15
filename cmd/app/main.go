package main

import (
	"TestCase/internal/configs"
	"TestCase/internal/controllers"
	"TestCase/internal/db"
	"TestCase/internal/redis"
	"TestCase/internal/repository"
	"TestCase/internal/routes"
	"TestCase/internal/services"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configs.InitRedis()
	configs.InitKafka()
	db.InitDatabase()

	router := gin.Default()

	personRepo := repository.NewPersonRepository(db.DB)
	redisClient := configs.InitRedis()
	enrichmentService := services.NewEnrichmentService(
		&http.Client{},
		&http.Client{},
		&http.Client{},
		personRepo,
	)

	cacheClient := redis.NewCacheClient(redisClient)
	apiController := controllers.NewAPIController(personRepo, *cacheClient, enrichmentService)

	apiRoutes := routes.NewRoutes(*apiController)
	apiRoutes.SetupAPIRoutes(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	fmt.Println("Server is running on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	}

	fmt.Println("Server gracefully stopped")
}
