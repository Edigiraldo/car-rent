package ports

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
)

type Database interface {
	Close() error
	InsertCar(ctx context.Context, dc domain.Car) (err error)
}
