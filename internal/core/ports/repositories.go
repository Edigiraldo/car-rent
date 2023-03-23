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
	FullUpdate(ctx context.Context, dc domain.Car) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, city string, from_car_id string, limit uint16) ([]domain.Car, error)
}
