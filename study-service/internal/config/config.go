package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load(path string) error {
	if _, err := os.Stat(path); err == nil {
		if err := godotenv.Load(path); err != nil {
			return err
		}
		log.Println("config: loaded from .env file")
		return nil
	}

	log.Println("config: .env file not found, using environment variables")
	return nil
}
