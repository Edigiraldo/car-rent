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

type UsersService interface {
	Register(ctx context.Context, car domain.User) (domain.User, error)
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	FullUpdate(ctx context.Context, du domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CitiesService interface {
	ListNames(ctx context.Context) ([]string, error)
}

type ReservationsService interface {
	Book(ctx context.Context, reservation domain.Reservation) (domain.Reservation, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Reservation, error)
	FullUpdate(ctx context.Context, dr domain.Reservation) error
}
