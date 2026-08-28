package main

import (
	"context"
	"log"
	"os"
	"time"

	"backend/internal/infrastructure/database"
	"backend/internal/infrastructure/firebase"
	"backend/internal/infrastructure/web"
	"backend/internal/interface/handler"
	"backend/internal/interface/repository"
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
	db := database.NewDB()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)

	// Usecases
	healthUsecase := usecase.NewHealthUsecase()
	authUsecase := usecase.NewAuthUsecase(fbClient, userRepo)
	eventUsecase := usecase.NewEventUsecase(eventRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler(healthUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	eventHandler := handler.NewEventHandler(eventUsecase)

	// Background Tasks: 終了時刻を過ぎたイベントを自動で finished に更新
	if db != nil {
		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			defer ticker.Stop()

			for range ticker.C {
				if err := eventUsecase.UpdateExpiredEvents(ctx); err != nil {
					log.Printf("error updating expired events: %v", err)
				}
			}
		}()
	} else {
		log.Println("Skipping expired event background task because database is unavailable")
	}

	// Router
	e := web.NewRouter(healthHandler, authHandler, eventHandler, fbClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
