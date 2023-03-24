package services

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/pkg/constants"
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
		return domain.Car{}, err
	}

	return dc, nil
}

func (c *Cars) FullUpdate(ctx context.Context, car domain.Car) error {
	return c.carsRepository.FullUpdate(ctx, car)
}

func (c *Cars) Delete(ctx context.Context, id uuid.UUID) error {
	return c.carsRepository.Delete(ctx, id)
}

// List cars by city name.
// from_car_id is the last document retrieved in the last page
func (c *Cars) List(ctx context.Context, city string, from_car_id string) ([]domain.Car, error) {
	if from_car_id == "" {
		from_car_id = constants.NULL_UUID
	}
	cars, err := c.carsRepository.List(ctx, city, from_car_id, constants.CARS_PER_PAGE)
	if err != nil {
		return []domain.Car{}, err
	}

	return cars, nil
}
