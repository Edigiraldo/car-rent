package main

import (
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers"
)

func BindRoutes(b *Server) {
	b.router.HandleFunc("/ping", handlers.Pong).Methods(http.MethodGet)
}
