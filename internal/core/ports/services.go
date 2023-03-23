package ports

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/google/uuid"
)

// mockgen -source=internal/core/ports/services.go -destination=internal/pkg/mocks/services.go

type CarsService interface {
	Register(ctx context.Context, car domain.Car) (domain.Car, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Car, error)
	FullUpdate(ctx context.Context, dc domain.Car) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, city string, from_car_id string) ([]domain.Car, error)
}
