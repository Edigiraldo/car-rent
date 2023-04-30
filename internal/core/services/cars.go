package services

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/pkg/constants"
	"github.com/google/uuid"
)

var (
	ErrCarNotFound     = "car not found"
	ErrCarNotAvailable = "car not available"
)

type Cars struct {
	carsRepository ports.CarsRepo
}

func NewCars(cr ports.CarsRepo) Cars {
	return Cars{
		carsRepository: cr,
	}
}

func (cs Cars) Register(ctx context.Context, car domain.Car) (domain.Car, error) {
	car.ID = uuid.New()

	if err := cs.carsRepository.Insert(ctx, car); err != nil {
		return domain.Car{}, err
	}

	return car, nil
}

func (cs Cars) Get(ctx context.Context, ID uuid.UUID) (domain.Car, error) {
	dc, err := cs.carsRepository.Get(ctx, ID)
	if err != nil {
		return domain.Car{}, err
	}

	return dc, nil
}

func (cs Cars) FullUpdate(ctx context.Context, car domain.Car) error {
	return cs.carsRepository.FullUpdate(ctx, car)
}

func (cs Cars) Delete(ctx context.Context, id uuid.UUID) error {
	return cs.carsRepository.Delete(ctx, id)
}

// List cars by city name.
// from_car_id is the last document retrieved in the last page
func (cs Cars) List(ctx context.Context, city string, from_car_id string) ([]domain.Car, error) {
	if from_car_id == "" {
		from_car_id = constants.Values.NULL_UUID
	}
	cars, err := cs.carsRepository.List(ctx, city, from_car_id, constants.Values.CARS_PER_PAGE)
	if err != nil {
		return []domain.Car{}, err
	}

	return cars, nil
}
