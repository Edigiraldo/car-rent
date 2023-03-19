package services

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/google/uuid"
)

type Cars struct {
	carsRepository ports.CarsRepo
}

func NewCars(cr ports.CarsRepo) *Cars {
	return &Cars{
		carsRepository: cr,
	}
}

func (c *Cars) Register(ctx context.Context, car domain.Car) (domain.Car, error) {
	car.ID = uuid.New()

	if err := c.carsRepository.InsertCar(ctx, car); err != nil {
		return domain.Car{}, err
	}

	return car, nil
}
