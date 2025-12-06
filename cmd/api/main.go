package main

import (
	"log"

	"github.com/Harschmann/hyper-checkout/internal/config"
	"github.com/Harschmann/hyper-checkout/internal/database"
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

	log.Printf("🚀 Server starting on port %s...", cfg.Port)

	
	select {}
}