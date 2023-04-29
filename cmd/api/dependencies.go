package main

import (
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/postgres"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers"
)

func initializeDependencies(config Config) (ports.Database, error) {
	// Initialize DB client
	carsRentDB, err := postgres.NewPostgresDB(config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Initialize repos
	citiesRepository := postgres.NewCitiesRepository(carsRentDB)
	carsRepository := postgres.NewCarsRepository(carsRentDB, citiesRepository)
	usersRepository := postgres.NewUsersRepository(carsRentDB)
	reservationsRepository := postgres.NewReservationsRepository(carsRentDB)

	// Initialize services
	carsService := services.NewCars(carsRepository, reservationsRepository)
	usersService := services.NewUsers(usersRepository, reservationsRepository)
	citiesService := services.NewCities(citiesRepository)
	reservationsService := services.NewReservations(reservationsRepository)

	//Initialize handlers
	healthHandler = handlers.NewHealth()
	carsHandler = handlers.NewCars(carsService)
	usersHandler = handlers.NewUsers(usersService)
	citiesHandler = handlers.NewCities(citiesService)
	reservationsHandler = handlers.NewReservations(reservationsService)

	return carsRentDB, nil
}
