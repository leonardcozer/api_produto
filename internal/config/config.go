package config

import (
	"os"
)

type Config struct {
	MongoURI string
	Database string
	Port     string
}

func Load() Config {
	return Config{
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Database: getEnv("MONGO_DB", "appdb"),
		Port:     getEnv("PORT", ":8080"),
	}
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}
