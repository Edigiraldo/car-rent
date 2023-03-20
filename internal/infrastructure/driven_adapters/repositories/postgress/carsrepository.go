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
