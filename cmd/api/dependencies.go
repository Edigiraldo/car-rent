package main

import (
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	postgres "github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/postgress"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers"
)

func initializeDependencies(config Config) (ports.Database, error) {
	// Initialize repos
	carsRentDB, err := postgres.NewPostgresDB(config.DatabaseURL)
	if err != nil {
		return nil, err
	}
	citiesRepository := postgres.NewCitiesRepository(carsRentDB)
	carsRepository := postgres.NewCarsRepository(carsRentDB, citiesRepository)
	usersRepository := postgres.NewUsersRepository(carsRentDB)

	// Initialize services
	carsService := services.NewCars(carsRepository)
	usersService := services.NewUsers(usersRepository)
	citiesService := services.NewCities(citiesRepository)

	//Initialize handlers
	healthHandler = handlers.NewHealth()
	carsHandler = handlers.NewCars(carsService)
	usersHandler = handlers.NewUsers(usersService)
	citiesHandler = handlers.NewCities(citiesService)

	return carsRentDB, nil
}
