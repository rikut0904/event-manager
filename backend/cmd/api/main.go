package main

import (
	"context"
	"log"
	"os"

	"backend/internal/infrastructure/firebase"
	"backend/internal/infrastructure/web"
	"backend/internal/interface/handler"
	"backend/internal/usecase"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ctx := context.Background()

	// Infrastructure
	fbClient := firebase.NewClient(ctx)

	// Usecases
	healthUsecase := usecase.NewHealthUsecase()
	authUsecase := usecase.NewAuthUsecase(fbClient)

	// Handlers
	healthHandler := handler.NewHealthHandler(healthUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)

	// Router
	e := web.NewRouter(healthHandler, authHandler, fbClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
