package postgres

import (
	"context"
	"database/sql"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/models"
)

type CarsRepo struct {
	db *sql.DB
}

func NewCarsRepository(db *sql.DB) *CarsRepo {
	return &CarsRepo{
		db: db,
	}
}

func (cr *CarsRepo) InsertCar(ctx context.Context, dc domain.Car) (err error) {
	car := models.LoadCarFromDomain(dc)
	_, err = cr.db.ExecContext(ctx, "INSERT INTO cars (id, type, seats, hourly_rent_cost, city, status) VALUES ($1, $2, $3, $4, $5, $6)",
		car.ID, car.Type, car.Seats, car.HourlyRentCost, car.City, car.Status)

	return err
}
