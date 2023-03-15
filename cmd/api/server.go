package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers"
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

func (b *Server) Config() *Config {
	return b.config
}

func (b *Server) BindRoutes() {
	b.router.HandleFunc("/ping", handlers.Pong).Methods(http.MethodGet)
}

func (b *Server) Start() {
	b.BindRoutes()
	log.Printf("Starting server on port %s\n", b.Config().Port)

	port := fmt.Sprintf(":%s", b.Config().Port)
	if err := http.ListenAndServe(port, b.router); err != nil {
		log.Fatal("Error while starting server: ", err)
	}
}
