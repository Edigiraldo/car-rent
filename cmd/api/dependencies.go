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
	carsRepository := postgres.NewCarsRepository(postgresDB.GetDBHandle())

	// Initialize services
	carsService := services.NewCars(carsRepository)

	//Initialize handlers
	healthHandler = handlers.NewHealth()
	carsHandler = handlers.NewCars(carsService)

	return nil
}
