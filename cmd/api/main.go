package main

import (
	"log"
	"net/http"

	"github.com/Harschmann/hyper-checkout/internal/config"
	"github.com/Harschmann/hyper-checkout/internal/database"
	"github.com/Harschmann/hyper-checkout/internal/handlers"
	"github.com/Harschmann/hyper-checkout/internal/repository"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewPostgresDB(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to instantiate database: %v", err)
	}
	defer db.Close()

	rdb := database.NewRedisClient(cfg.RedisUrl)
	defer rdb.Close()

	// 3. Setup Layers
	// Repository (Logic) -> Handler (HTTP)
	repo := repository.NewStore(db)
	handler := handlers.NewHandler(repo)

	// 4. Setup Router (Chi)
	r := chi.NewRouter()
	r.Use(middleware.Logger)    // Log every request
	r.Use(middleware.Recoverer) // Recover from panics
	// r.Use(myMiddleware.RateLimit(rdb))  // Block spammers before they hit the DB

	// 5. Define Routes
	r.Post("/purchase", handler.HandlePurchase)

	// 6. Start Server
	log.Printf("🚀 Server starting on port %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
