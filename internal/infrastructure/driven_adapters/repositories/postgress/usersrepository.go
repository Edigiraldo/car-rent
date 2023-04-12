package postgres

import (
	"context"
	"database/sql"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/models"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (cr *UsersRepo) Insert(ctx context.Context, du domain.User) (err error) {
	car := models.LoadUserFromDomain(du)
	_, err = cr.db.ExecContext(ctx, "INSERT INTO users (id, first_name, last_name, email, type, status) VALUES ($1, $2, $3, $4, $5, $6)",
		car.ID, car.FirstName, car.LastName, car.Email, car.Type, car.Status)

	return err
}
