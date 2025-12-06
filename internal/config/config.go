package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl    string
	RedisUrl string
	Port     string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env vars")
	}
	return Config{
		DBUrl:    getEnv("DB_URL", "postgres://user:password@localhost:5432/hypercheckout?sslmode=disable"),
		RedisUrl: getEnv("REDIS_URL", "localhost:6379"),
		Port:     getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
