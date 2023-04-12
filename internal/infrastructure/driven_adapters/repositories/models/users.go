package models

import (
	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
}

func (u *User) ToDomain() domain.User {
	return domain.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Type:      u.Type,
		Status:    u.Status,
	}
}

func LoadUserFromDomain(du domain.User) User {
	return User{
		ID:        du.ID,
		FirstName: du.FirstName,
		LastName:  du.LastName,
		Email:     du.Email,
		Type:      du.Type,
		Status:    du.Status,
	}
}
