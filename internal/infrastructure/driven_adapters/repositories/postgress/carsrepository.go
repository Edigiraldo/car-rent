package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/models"
	"github.com/google/uuid"
)

type CarsRepo struct {
	db *sql.DB
}

func NewCarsRepository(db *sql.DB) *CarsRepo {
	return &CarsRepo{
		db: db,
	}
}

func (cr *CarsRepo) Insert(ctx context.Context, dc domain.Car) (err error) {
	car := models.LoadCarFromDomain(dc)
	_, err = cr.db.ExecContext(ctx, "INSERT INTO cars (id, type, seats, hourly_rent_cost, city, status) VALUES ($1, $2, $3, $4, $5, $6)",
		car.ID, car.Type, car.Seats, car.HourlyRentCost, car.City, car.Status)

	return err
}

func (cr *CarsRepo) Get(ctx context.Context, ID uuid.UUID) (dc domain.Car, err error) {
	var car models.Car
	if err := cr.db.QueryRowContext(ctx, "SELECT * FROM cars WHERE ID = $1", ID).
		Scan(&car.ID, &car.Type, &car.Seats, &car.HourlyRentCost, &car.City, &car.Status); err != nil {
		if err == sql.ErrNoRows {
			return domain.Car{}, errors.New(services.ErrCarNotFound)
		}
		return domain.Car{}, err
	}

	return car.ToDomain(), nil
}

// Updates car row. If car was not found returns an error.
func (cr *CarsRepo) FullUpdate(ctx context.Context, dc domain.Car) error {
	car := models.LoadCarFromDomain(dc)

	result, err := cr.db.ExecContext(ctx, "UPDATE cars SET type=$1, seats=$2, hourly_rent_cost=$3, city=$4, status=$5 WHERE id=$6",
		car.Type, car.Seats, car.HourlyRentCost, car.City, car.Status, car.ID)
	if err != nil {
		return err
	}

	numUpdatedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numUpdatedRows == 0 {
		return errors.New(services.ErrCarNotFound)
	}

	return nil
}

func (cr *CarsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := cr.db.ExecContext(ctx, "Delete FROM cars WHERE id=$1", id)

	return err
}
