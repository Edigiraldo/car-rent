package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("PORT environment variable was nos found")
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		log.Fatal("DATABASE_URL environment variable was nos found")
	}

	config := Config{
		Port:        PORT,
		DatabaseURL: DATABASE_URL,
	}
	s, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	s.Start()
}
