package main

import (
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/ports"
)

var (
	healthHandler ports.HeathController
	carsHandler   ports.CarsController
)

func BindRoutes(b *Server) {
	// Health routes
	b.router.HandleFunc("/api/v1/ping", healthHandler.Pong).Methods(http.MethodGet)

	// Car routes
	b.router.HandleFunc("/api/v1/cars", carsHandler.Register).Methods(http.MethodPost)
	b.router.HandleFunc("/api/v1/cars/{id}", carsHandler.Get).Methods(http.MethodGet)
	b.router.HandleFunc("/api/v1/cars/{id}", carsHandler.FullUpdate).Methods(http.MethodPut)
	b.router.HandleFunc("/api/v1/cars/{id}", carsHandler.Delete).Methods(http.MethodDelete)
}
