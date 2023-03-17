package main

import (
	"github.com/Edigiraldo/car-rent/internal/core/services"
	postgres "github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/postgress"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers"
)

func initializeDependencies(config Config) error {
	carsRepository, err := postgres.NewPostgresDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	carsService := services.NewCars(carsRepository)

	healthHandler = handlers.NewHealth()
	carsHandler = handlers.NewCars(carsService)

	return nil
}
