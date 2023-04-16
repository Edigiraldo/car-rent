package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/models"
	"github.com/google/uuid"
)

type CarsRepo struct {
	db               *sql.DB
	citiesRepository ports.CitiesRepo
}

func NewCarsRepository(db *sql.DB, cr ports.CitiesRepo) *CarsRepo {
	return &CarsRepo{
		db:               db,
		citiesRepository: cr,
	}
}

func (cr *CarsRepo) Insert(ctx context.Context, dc domain.Car) (err error) {
	car := models.LoadCarFromDomain(dc)

	car.CityID, err = cr.citiesRepository.GetIdByName(ctx, dc.CityName)
	if err != nil {
		return err
	}

	_, err = cr.db.ExecContext(ctx, "INSERT INTO cars (id, type, seats, hourly_rent_cost, city_id, status) VALUES ($1, $2, $3, $4, $5, $6)",
		car.ID, car.Type, car.Seats, car.HourlyRentCost, car.CityID, car.Status)

	return err
}

func (cr *CarsRepo) Get(ctx context.Context, ID uuid.UUID) (dc domain.Car, err error) {
	var car models.Car
	if err := cr.db.QueryRowContext(ctx, "SELECT * FROM cars WHERE ID = $1", ID).
		Scan(&car.ID, &car.Type, &car.Seats, &car.HourlyRentCost, &car.CityID, &car.Status); err != nil {
		if err == sql.ErrNoRows {
			return domain.Car{}, errors.New(services.ErrCarNotFound)
		}
		return domain.Car{}, err
	}

	cityName, err := cr.citiesRepository.GetNameByID(ctx, car.CityID)
	if err != nil {
		return domain.Car{}, err
	}

	return car.ToDomain(cityName), nil
}

// Updates car row. If car was not found returns an error.
func (cr *CarsRepo) FullUpdate(ctx context.Context, dc domain.Car) (err error) {
	car := models.LoadCarFromDomain(dc)

	if car.CityID, err = cr.citiesRepository.GetIdByName(ctx, dc.CityName); err != nil {
		return err
	}

	result, err := cr.db.ExecContext(ctx, "UPDATE cars SET type=$1, seats=$2, hourly_rent_cost=$3, city_id=$4, status=$5 WHERE id=$6",
		car.Type, car.Seats, car.HourlyRentCost, car.CityID, car.Status, car.ID)
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
	_, err := cr.db.ExecContext(ctx, "DELETE FROM cars WHERE id=$1", id)

	return err
}

// List cars by city name.
// from_car_id is the last document retrieved in the last page.
// limit is the number of documents per page.
func (cr *CarsRepo) List(ctx context.Context, cityName string, from_car_id string, limit uint16) ([]domain.Car, error) {
	var cars []domain.Car

	cityID, err := cr.citiesRepository.GetIdByName(ctx, cityName)
	if err != nil {
		return nil, err
	}

	rows, err := cr.db.QueryContext(ctx, "SELECT * FROM cars WHERE city_id=$1 AND id > $2 ORDER BY id ASC LIMIT $3", cityID, from_car_id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		car := models.Car{}
		if err := rows.Scan(&car.ID, &car.Type, &car.Seats, &car.HourlyRentCost, &car.CityID, &car.Status); err != nil {
			return nil, err
		}

		cars = append(cars, car.ToDomain(cityName))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}
