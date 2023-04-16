package services

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/ports"
)

var (
	ErrInvalidCityName = "city name is not valid"
	ErrCityNotFound    = "city not found"
)

type Cities struct {
	carsRepository ports.CitiesRepo
}

func NewCities(cr ports.CitiesRepo) *Cities {
	return &Cities{
		carsRepository: cr,
	}
}

func (cs *Cities) ListNames(ctx context.Context) ([]string, error) {
	return cs.carsRepository.ListNames(ctx)
}
