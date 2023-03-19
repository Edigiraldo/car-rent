package services

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
)

type Cars struct {
	carsRepository ports.CarsRepo
}

func NewCars(cr ports.CarsRepo) *Cars {
	return &Cars{
		carsRepository: cr,
	}
}

func (c *Cars) Register(ctx context.Context, car domain.Car) error {
	return c.carsRepository.InsertCar(ctx, car)
}
