package ports

import (
	"context"
	"database/sql"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/google/uuid"
)

// mockgen -source=internal/core/ports/repositories.go -destination=internal/pkg/mocks/repositories.go

type Database interface {
	GetDBHandle() *sql.DB
	Close() error
}

type CarsRepo interface {
	Insert(ctx context.Context, dc domain.Car) (err error)
	Get(ctx context.Context, ID uuid.UUID) (dc domain.Car, err error)
	FullUpdate(ctx context.Context, dc domain.Car) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, cityName string, from_car_id string, limit uint16) ([]domain.Car, error)
}

type UsersRepo interface {
	Insert(ctx context.Context, du domain.User) (err error)
	Get(ctx context.Context, ID uuid.UUID) (dc domain.User, err error)
	FullUpdate(ctx context.Context, du domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CitiesRepo interface {
	GetIdByName(ctx context.Context, name string) (ID uuid.UUID, err error)
	GetNameByID(ctx context.Context, ID uuid.UUID) (name string, err error)
	ListNames(ctx context.Context) ([]string, error)
}

type ReservationsRepo interface {
	Insert(ctx context.Context, dr domain.Reservation) (err error)
}
