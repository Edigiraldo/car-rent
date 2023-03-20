package services

import (
	"context"
	"errors"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrCarNotFound = "car not found"
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

	if err := c.carsRepository.Insert(ctx, car); err != nil {
		return domain.Car{}, err
	}

	return car, nil
}

func (c *Cars) Get(ctx context.Context, ID uuid.UUID) (domain.Car, error) {
	dc, err := c.carsRepository.Get(ctx, ID)
	if err != nil {
		if err.Error() == ErrCarNotFound {
			return domain.Car{}, errors.New(ErrCarNotFound)
		}
		return domain.Car{}, err
	}

	return dc, nil
}
