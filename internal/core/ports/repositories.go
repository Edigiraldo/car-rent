package ports

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
)

// mockgen -source=internal/core/ports/repositories.go -destination=internal/pkg/mocks/repositories.go

type Database interface {
	Close() error
	InsertCar(ctx context.Context, dc domain.Car) (err error)
}
