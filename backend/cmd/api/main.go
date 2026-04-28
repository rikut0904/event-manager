package main

import (
	"backend/internal/infrastructure/web"
	"backend/internal/interface/handler"
	"backend/internal/usecase"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Dependency Injection
	healthUseCase := usecase.NewHealthUseCase()
	healthHandler := handler.NewHealthHandler(healthUseCase)
	mux := web.NewRouter(healthHandler)

	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
