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

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (cr *UsersRepo) Insert(ctx context.Context, du domain.User) (err error) {
	user := models.LoadUserFromDomain(du)
	_, err = cr.db.ExecContext(ctx, "INSERT INTO users (id, first_name, last_name, email, type, status) VALUES ($1, $2, $3, $4, $5, $6)",
		user.ID, user.FirstName, user.LastName, user.Email, user.Type, user.Status)

	return err
}

func (cr *UsersRepo) Get(ctx context.Context, ID uuid.UUID) (dc domain.User, err error) {
	var user models.User
	if err := cr.db.QueryRowContext(ctx, "SELECT * FROM users WHERE ID = $1", ID).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Type, &user.Status); err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errors.New(services.ErrUserNotFound)
		}
		return domain.User{}, err
	}

	return user.ToDomain(), nil
}
