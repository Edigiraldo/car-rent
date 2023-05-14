package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UsersRepo struct {
	ports.Database
}

func NewUsersRepository(db ports.Database) *UsersRepo {
	return &UsersRepo{
		Database: db,
	}
}

func (ur *UsersRepo) Insert(ctx context.Context, du domain.User) (err error) {
	user := models.LoadUserFromDomain(du)
	_, err = ur.GetDBHandle().ExecContext(ctx, "INSERT INTO users (id, first_name, last_name, email, type, status) VALUES ($1, $2, $3, $4, $5, $6)",
		user.ID, user.FirstName, user.LastName, user.Email, user.Type, user.Status)

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		if strings.Contains(pqErr.Message, "unique_email") {
			return errors.New(services.ErrEmailAlreadyRegistered)
		}
	}

	return err
}

func (ur *UsersRepo) Get(ctx context.Context, ID uuid.UUID) (dc domain.User, err error) {
	var user models.User
	if err := ur.GetDBHandle().QueryRowContext(ctx, "SELECT * FROM users WHERE ID = $1", ID).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Type, &user.Status); err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errors.New(services.ErrUserNotFound)
		}
		return domain.User{}, err
	}

	return user.ToDomain(), nil
}

func (ur *UsersRepo) FullUpdate(ctx context.Context, dc domain.User) error {
	user := models.LoadUserFromDomain(dc)

	result, err := ur.GetDBHandle().ExecContext(ctx, "UPDATE users SET first_name=$1, last_name=$2, email=$3, type=$4, status=$5 WHERE id=$6",
		user.FirstName, user.LastName, user.Email, user.Type, user.Status, user.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Message, "unique_email") {
				return errors.New(services.ErrEmailAlreadyRegistered)
			}
		}
		return err
	}

	numUpdatedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numUpdatedRows == 0 {
		return errors.New(services.ErrUserNotFound)
	}

	return nil
}

func (ur *UsersRepo) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := ur.GetDBHandle().ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}

	numDeletedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numDeletedRows == 0 {
		return errors.New(services.ErrUserNotFound)
	}

	return err
}
