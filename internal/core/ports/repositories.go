package ports

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/google/uuid"
)

// mockgen -source=internal/core/ports/repositories.go -destination=internal/pkg/mocks/repositories.go

type Database interface {
	GetDBHandle()
	Close() error
}

type CarsRepo interface {
	Insert(ctx context.Context, dc domain.Car) (err error)
	Get(ctx context.Context, ID uuid.UUID) (dc domain.Car, err error)
}
