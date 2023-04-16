package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/pkg/constants"
	"github.com/gorilla/mux"
)

var (
	ErrEmptyDabaseURL = "database url must be specified in server configuration"
	ErrEmptyPort      = "port must be specified in server configuration"
)

type Config struct {
	Port        string
	DatabaseURL string
}

type Server struct {
	config *Config
	router *mux.Router
}

func NewServer(config Config) (*Server, error) {
	if config.Port == "" {
		return nil, errors.New(ErrEmptyPort)
	}

	if config.DatabaseURL == "" {
		return nil, errors.New(ErrEmptyDabaseURL)
	}

	server := &Server{
		config: &config,
		router: mux.NewRouter(),
	}

	return server, nil
}

func (b *Server) Start() {
	carsRentDB, err := initializeDependencies(*b.config)
	if err != nil {
		log.Fatal("error while initializing dependencies: ", err)
	}
	defer carsRentDB.Close()

	if err := constants.InitValues(); err != nil {
		log.Fatal("error while loading constants file: ", err)
	}

	BindRoutes(b)
	log.Printf("starting server on port %s\n", b.config.Port)

	port := fmt.Sprintf(":%s", b.config.Port)
	if err := http.ListenAndServe(port, b.router); err != nil {
		log.Fatal("error while starting server: ", err)
	}
}
