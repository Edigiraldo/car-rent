package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/google/uuid"
)

type CitiesRepo struct {
	ports.Database
}

func NewCitiesRepository(db ports.Database) *CitiesRepo {
	return &CitiesRepo{
		Database: db,
	}
}

func (cr *CitiesRepo) GetIdByName(ctx context.Context, name string) (ID uuid.UUID, err error) {
	if err := cr.GetDBHandle().QueryRowContext(ctx, "SELECT id FROM cities WHERE name = $1", name).
		Scan(&ID); err != nil {
		if err == sql.ErrNoRows {
			return uuid.UUID{}, errors.New(services.ErrInvalidCityName)
		}
		return uuid.UUID{}, err
	}

	return ID, nil
}

func (cr *CitiesRepo) GetNameByID(ctx context.Context, ID uuid.UUID) (name string, err error) {
	if err := cr.GetDBHandle().QueryRowContext(ctx, "SELECT name FROM cities WHERE id = $1", ID).
		Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New(services.ErrCityNotFound)
		}
		return "", err
	}

	return name, nil
}

func (cr *CitiesRepo) ListNames(ctx context.Context) ([]string, error) {
	var cityNames []string

	rows, err := cr.GetDBHandle().QueryContext(ctx, "SELECT name FROM cities ORDER BY name LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var cityName string
		if err := rows.Scan(&cityName); err != nil {
			return nil, err
		}

		cityNames = append(cityNames, cityName)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cityNames, nil
}
