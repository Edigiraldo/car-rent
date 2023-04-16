package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/google/uuid"
)

type CitiesRepo struct {
	db *sql.DB
}

func NewCitiesRepository(db *sql.DB) *CitiesRepo {
	return &CitiesRepo{
		db: db,
	}
}

func (cr *CitiesRepo) GetIdByName(ctx context.Context, name string) (ID uuid.UUID, err error) {
	if err := cr.db.QueryRowContext(ctx, "SELECT id FROM cities WHERE name = $1", name).
		Scan(&ID); err != nil {
		if err == sql.ErrNoRows {
			return uuid.UUID{}, errors.New(services.ErrInvalidCityName)
		}
		return uuid.UUID{}, err
	}

	return ID, nil
}

func (cr *CitiesRepo) GetNameByID(ctx context.Context, ID uuid.UUID) (name string, err error) {
	if err := cr.db.QueryRowContext(ctx, "SELECT name FROM cities WHERE id = $1", ID).
		Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New(services.ErrCityNotFound)
		}
		return "", err
	}

	return name, nil
}
