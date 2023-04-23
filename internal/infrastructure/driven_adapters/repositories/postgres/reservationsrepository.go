package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/models"
	"github.com/lib/pq"
)

type ReservationsRepo struct {
	ports.Database
}

func NewReservationsRepository(db ports.Database) *ReservationsRepo {
	return &ReservationsRepo{
		Database: db,
	}
}

func (cr *ReservationsRepo) Insert(ctx context.Context, dc domain.Reservation) (err error) {
	reservation := models.LoadReservationFromDomain(dc)

	_, err = cr.GetDBHandle().ExecContext(ctx, "INSERT INTO reservations (id, user_id, car_id, status, payment_status, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		reservation.ID, reservation.UserID, reservation.CarID, reservation.Status, reservation.PaymentStatus, reservation.StartDate, reservation.EndDate)

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
		if strings.Contains(pqErr.Message, "user_id") {
			return errors.New(services.ErrUserNotFound)
		} else if strings.Contains(pqErr.Message, "car_id") {
			return errors.New(services.ErrCarNotFound)
		}
	}

	return err
}
