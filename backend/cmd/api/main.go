package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"backend/internal/infrastructure/commonid"
	"backend/internal/infrastructure/database"
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
	appOrigin, err := requiredEnv("APP_ORIGIN")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Infrastructure
	db := database.NewDB()
	commonID, err := commonid.New(commonid.Config{
		Origin:            envOrDefault("COMMON_ID_ORIGIN", "http://localhost:13000"),
		APIOrigin:         envOrDefault("COMMON_ID_API_ORIGIN", "http://localhost:18080"),
		ClientID:          os.Getenv("COMMON_ID_CLIENT_ID"),
		RedirectURI:       envOrDefault("COMMON_ID_REDIRECT_URI", "http://localhost:8080/auth/callback"),
		LogoutRedirectURI: envOrDefault("COMMON_ID_LOGOUT_REDIRECT_URI", "http://localhost:8080/auth/logout/callback"),
		APIKey:            os.Getenv("COMMON_ID_API_KEY"),
	})
	if err != nil {
		log.Printf("WARNING: Common ID is not configured: %v", err)
		commonID = nil
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)

	// Usecases
	healthUsecase := usecase.NewHealthUsecase()
	authUsecase := usecase.NewAuthUsecase(nil, userRepo)
	eventUsecase := usecase.NewEventUsecase(eventRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler(healthUsecase)
	authHandler := handler.NewAuthHandler(authUsecase, commonID, appOrigin)
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
	e := web.NewRouter(healthHandler, authHandler, eventHandler, commonID, appOrigin)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func requiredEnv(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("%s is required", name)
	}
	return value, nil
}
