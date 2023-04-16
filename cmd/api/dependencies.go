package main

import (
	"github.com/Edigiraldo/car-rent/internal/core/services"
	postgres "github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/postgress"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers"
)

func initializeDependencies(config Config) error {
	// Initialize repos
	postgresDB, err := postgres.NewPostgresDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	citiesRepository := postgres.NewCitiesRepository(postgresDB.GetDBHandle())
	carsRepository := postgres.NewCarsRepository(postgresDB.GetDBHandle(), citiesRepository)
	usersRepository := postgres.NewUsersRepository(postgresDB.GetDBHandle())

	// Initialize services
	carsService := services.NewCars(carsRepository)
	usersService := services.NewUsers(usersRepository)

	//Initialize handlers
	healthHandler = handlers.NewHealth()
	carsHandler = handlers.NewCars(carsService)
	usersHandler = handlers.NewUsers(usersService)

	return nil
}
