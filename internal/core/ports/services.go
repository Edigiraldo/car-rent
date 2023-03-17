package ports

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
)

type CarsService interface {
	Register(ctx context.Context, car domain.Car) error
}
