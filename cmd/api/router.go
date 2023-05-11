package main

import (
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/ports"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	healthHandler       ports.HeathController
	carsHandler         ports.CarsController
	usersHandler        ports.UsersController
	citiesHandler       ports.CitiesController
	reservationsHandler ports.ReservationsController
)

func BindRoutes(b *Server) {
	// Swagger
	b.router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// Health routes
	b.router.HandleFunc("/api/v1/ping", healthHandler.Pong).Methods(http.MethodGet)

	// Cars routes
	b.router.HandleFunc("/api/v1/cars", carsHandler.Register).Methods(http.MethodPost)
	b.router.HandleFunc("/api/v1/cars/{id}", carsHandler.Get).Methods(http.MethodGet)
	b.router.HandleFunc("/api/v1/cars/{id}", carsHandler.FullUpdate).Methods(http.MethodPut)
	b.router.HandleFunc("/api/v1/cars/{id}", carsHandler.Delete).Methods(http.MethodDelete)
	b.router.HandleFunc("/api/v1/cars/", carsHandler.List).Methods(http.MethodGet)
	b.router.HandleFunc("/api/v1/cars/{id}/reservations", reservationsHandler.GetByCarID).Methods(http.MethodGet)

	// Users routes
	b.router.HandleFunc("/api/v1/users", usersHandler.SignUp).Methods(http.MethodPost)
	b.router.HandleFunc("/api/v1/users/{id}", usersHandler.Get).Methods(http.MethodGet)
	b.router.HandleFunc("/api/v1/users/{id}", usersHandler.FullUpdate).Methods(http.MethodPut)
	b.router.HandleFunc("/api/v1/users/{id}", usersHandler.Delete).Methods(http.MethodDelete)
	b.router.HandleFunc("/api/v1/users/{id}/reservations", reservationsHandler.GetByUserID).Methods(http.MethodGet)

	// Cities routes
	b.router.HandleFunc("/api/v1/cities/names", citiesHandler.ListNames).Methods(http.MethodGet)

	// Reservations routes
	b.router.HandleFunc("/api/v1/reservations", reservationsHandler.Book).Methods(http.MethodPost)
	b.router.HandleFunc("/api/v1/reservations/{id}", reservationsHandler.Get).Methods(http.MethodGet)
	b.router.HandleFunc("/api/v1/reservations/{id}", reservationsHandler.FullUpdate).Methods(http.MethodPut)
	b.router.HandleFunc("/api/v1/reservations/{id}", reservationsHandler.Delete).Methods(http.MethodDelete)
}
