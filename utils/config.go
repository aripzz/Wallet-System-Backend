package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetJWTSecret() string {
	LoadEnv()
	return os.Getenv("JWT_SECRET")
}
