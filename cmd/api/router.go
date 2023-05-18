package main

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers"
	"github.com/Edigiraldo/car-rent/pkg/httphandler"
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
	// Recovery middleware
	b.router.Use(recovery)

	// Use strict slashes
	b.router.StrictSlash(true)

	rv1 := b.router.PathPrefix("/api/v1").Subrouter()

	// Swagger
	rv1.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// Health routes
	rv1.HandleFunc("/ping", healthHandler.Pong).Methods(http.MethodGet)

	// Cars routes
	rv1.HandleFunc("/cars", carsHandler.Register).Methods(http.MethodPost)
	rv1.HandleFunc("/cars/{id}", carsHandler.Get).Methods(http.MethodGet)
	rv1.HandleFunc("/cars/{id}", carsHandler.FullUpdate).Methods(http.MethodPut)
	rv1.HandleFunc("/cars/{id}", carsHandler.Delete).Methods(http.MethodDelete)
	rv1.HandleFunc("/cars/", carsHandler.List).Methods(http.MethodGet)

	// Users routes
	rv1.HandleFunc("/users", usersHandler.SignUp).Methods(http.MethodPost)
	rv1.HandleFunc("/users/{id}", usersHandler.Get).Methods(http.MethodGet)
	rv1.HandleFunc("/users/{id}", usersHandler.FullUpdate).Methods(http.MethodPut)
	rv1.HandleFunc("/users/{id}", usersHandler.Delete).Methods(http.MethodDelete)

	// Cities routes
	rv1.HandleFunc("/cities/names", citiesHandler.ListNames).Methods(http.MethodGet)

	// Reservations routes
	rv1.HandleFunc("/reservations", reservationsHandler.Book).Methods(http.MethodPost)
	rv1.HandleFunc("/reservations/{id}", reservationsHandler.Get).Methods(http.MethodGet)
	rv1.HandleFunc("/reservations/{id}", reservationsHandler.FullUpdate).Methods(http.MethodPut)
	rv1.HandleFunc("/reservations/{id}", reservationsHandler.Delete).Methods(http.MethodDelete)
	rv1.HandleFunc("/reservations", reservationsHandler.List).Methods(http.MethodGet)
	rv1.HandleFunc("/cars/{id}/reservations", reservationsHandler.GetByCarID).Methods(http.MethodGet)
	rv1.HandleFunc("/users/{id}/reservations", reservationsHandler.GetByUserID).Methods(http.MethodGet)

}

func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic occurred: %v", err)
				debug.PrintStack()
				httphandler.WriteErrorResponse(w, http.StatusInternalServerError, handlers.ErrInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
