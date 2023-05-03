package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/models"
	"github.com/google/uuid"
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

func (rr ReservationsRepo) Insert(ctx context.Context, dc domain.Reservation) (err error) {
	reservation := models.LoadReservationFromDomain(dc)

	_, err = rr.GetDBHandle().ExecContext(ctx, "INSERT INTO reservations (id, user_id, car_id, status, payment_status, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6, $7)",
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

func (rr ReservationsRepo) Get(ctx context.Context, ID uuid.UUID) (dc domain.Reservation, err error) {
	var reservation models.Reservation
	if err := rr.GetDBHandle().QueryRowContext(ctx, "SELECT * FROM reservations WHERE ID = $1", ID).
		Scan(&reservation.ID, &reservation.UserID, &reservation.CarID, &reservation.Status, &reservation.PaymentStatus, &reservation.StartDate, &reservation.EndDate); err != nil {
		if err == sql.ErrNoRows {
			return domain.Reservation{}, errors.New(services.ErrReservationNotFound)
		}
		return domain.Reservation{}, err
	}

	return reservation.ToDomain(), nil
}

func (rr ReservationsRepo) FullUpdate(ctx context.Context, dr domain.Reservation) (err error) {
	reservation := models.LoadReservationFromDomain(dr)

	result, err := rr.GetDBHandle().ExecContext(ctx, "UPDATE reservations SET user_id=$1, car_id=$2, status=$3, payment_status=$4, start_date=$5, end_date=$6 WHERE id=$7",
		reservation.UserID, reservation.CarID, reservation.Status, reservation.PaymentStatus, reservation.StartDate, reservation.EndDate, reservation.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			if strings.Contains(pqErr.Message, "user_id") {
				return errors.New(services.ErrUserNotFound)
			} else if strings.Contains(pqErr.Message, "car_id") {
				return errors.New(services.ErrCarNotFound)
			}
		}

		return err
	}

	numUpdatedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numUpdatedRows == 0 {
		return errors.New(services.ErrReservationNotFound)
	}

	return nil
}

func (rr ReservationsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := rr.GetDBHandle().ExecContext(ctx, "DELETE FROM reservations WHERE id=$1", id)

	return err
}

func (rr ReservationsRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (dr []domain.Reservation, err error) {
	var reservations []domain.Reservation

	rows, err := rr.GetDBHandle().QueryContext(ctx, "SELECT * FROM reservations WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		reservation := models.Reservation{}
		if err := rows.Scan(&reservation.ID, &reservation.UserID, &reservation.CarID, &reservation.Status, &reservation.PaymentStatus, &reservation.StartDate, &reservation.EndDate); err != nil {
			return nil, err
		}

		reservations = append(reservations, reservation.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

func (rr ReservationsRepo) GetByCarID(ctx context.Context, carID uuid.UUID) (dr []domain.Reservation, err error) {
	var reservations []domain.Reservation

	rows, err := rr.GetDBHandle().QueryContext(ctx, "SELECT * FROM reservations WHERE car_id=$1", carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		reservation := models.Reservation{}
		if err := rows.Scan(&reservation.ID, &reservation.UserID, &reservation.CarID, &reservation.Status, &reservation.PaymentStatus, &reservation.StartDate, &reservation.EndDate); err != nil {
			return nil, err
		}

		reservations = append(reservations, reservation.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

func (rr ReservationsRepo) GetByCarIDAndTimeFrame(ctx context.Context, carID uuid.UUID, startDate time.Time, endDate time.Time) (dr []domain.Reservation, err error) {
	var reservations []domain.Reservation

	query := "SELECT * FROM reservations WHERE car_id=$1 AND start_date BETWEEN $2 AND $3 AND end_date BETWEEN $2 AND $3"
	rows, err := rr.GetDBHandle().QueryContext(ctx, query, carID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		reservation := models.Reservation{}
		if err := rows.Scan(&reservation.ID, &reservation.UserID, &reservation.CarID, &reservation.Status, &reservation.PaymentStatus, &reservation.StartDate, &reservation.EndDate); err != nil {
			return nil, err
		}

		reservations = append(reservations, reservation.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}
