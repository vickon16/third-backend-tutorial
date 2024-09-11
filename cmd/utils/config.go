package utils

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URL                 string
	PORT                   string
	JWT_SECRET             string
	JWTExpirationInSeconds int64
}

// create a singleton
var Configs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		DB_URL:                 getEnv("DB_URL", "root:Password123!@tcp(localhost:3306)/go_test_db?parseTime=true"),
		PORT:                   getEnv("PORT", "4000"),
		JWT_SECRET:             getEnv("JWT_SECRET", "openSecret"),
		JWTExpirationInSeconds: getEnvInt("JWT_EXP", 60*60 /* 1 hour */),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return intValue
	}

	return fallback
}
