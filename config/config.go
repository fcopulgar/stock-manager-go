package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("The .env file could not be loaded, system environment variables will be used.")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
