package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/Edigiraldo/car-rent/doc/swagger"
	"github.com/joho/godotenv"
	"github.com/swaggo/swag"
)

// @title Car Rent API
// @description This is an API to manage a car rent service
// @version 1.0
// @host localhost:5050
// @BasePath /api/v1/
func main() {

	envPath := getEnvPath()

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal(err)
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("PORT environment variable was not found")
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		log.Fatal("DATABASE_URL environment variable was not found")
	}

	swag.SetCodeExampleFilesDirectory("../../doc")

	config := Config{
		Port:        PORT,
		DatabaseURL: DATABASE_URL,
	}

	fmt.Println("Config:", config)

	s, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	s.Start()
}

func getEnvPath() string {
	ENVIRONMENT := os.Getenv("ENVIRONMENT")
	switch ENVIRONMENT {
	case "local":
		return ".env.local"
	case "debug":
		return ".env.debug"
	default:
		return ".env"
	}

}
