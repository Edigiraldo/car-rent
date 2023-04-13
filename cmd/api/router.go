package main

import (
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/ports"
)

var (
	healthHandler ports.HeathController
	carsHandler   ports.CarsController
	usersHandler  ports.UsersController
)

func BindRoutes(b *Server) {
	// Health routes
	b.router.HandleFunc("/api/v1/ping", healthHandler.Pong).Methods(http.MethodGet)

	// Cars routes
	b.router.HandleFunc("/api/v1/cars", carsHandler.Register).Methods(http.MethodPost)
	b.router.HandleFunc("/api/v1/cars/{id}", carsHandler.Get).Methods(http.MethodGet)
	b.router.HandleFunc("/api/v1/cars/{id}", carsHandler.FullUpdate).Methods(http.MethodPut)
	b.router.HandleFunc("/api/v1/cars/{id}", carsHandler.Delete).Methods(http.MethodDelete)
	b.router.HandleFunc("/api/v1/cars/", carsHandler.List).Methods(http.MethodGet)

	// Users routes
	b.router.HandleFunc("/api/v1/users", usersHandler.SignUp).Methods(http.MethodPost)
	b.router.HandleFunc("/api/v1/users/{id}", usersHandler.Get).Methods(http.MethodGet)

}
