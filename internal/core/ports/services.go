package ports

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
)

//mockgen -source=internal/core/ports/services.go -destination=internal/pkg/mocks/services.go

type CarsService interface {
	Register(ctx context.Context, car domain.Car) error
}
