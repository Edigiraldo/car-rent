package dtos

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/pkg/constants"
	"github.com/Edigiraldo/car-rent/pkg/utils"
	"github.com/google/uuid"
)

var (
	ErrEmptyFirstName    = "first name cannot be empty"
	ErrEmptyLastName     = "last name cannot be empty"
	ErrInvalidEmail      = "invalid email"
	ErrInvalidUserType   = "invalid user type"
	ErrInvalidUserStatus = "invalid user status"
)

type ListUsersResponse struct {
	Users []User `json:"users"`
}

type User struct {
	ID        uuid.UUID `json:"id,omitempty"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
}

func (u User) ToDomain() domain.User {
	return domain.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Type:      u.Type,
		Status:    u.Status,
	}
}

func (u *User) FromDomain(du domain.User) {
	u.ID = du.ID
	u.FirstName = du.FirstName
	u.LastName = du.LastName
	u.Email = du.Email
	u.Type = du.Type
	u.Status = du.Status
}

func UserFromBody(body io.Reader) (User, error) {
	var user User
	err := json.NewDecoder(body).Decode(&user)
	if err != nil {
		return User{}, err
	}

	if user.FirstName == "" {
		return User{}, errors.New(ErrEmptyFirstName)
	}

	if user.LastName == "" {
		return User{}, errors.New(ErrEmptyLastName)
	}

	if !isValidEmail(user.Email) {
		return User{}, errors.New(ErrInvalidEmail)
	}

	if !isValidUserType(user.Type) {
		return User{}, errors.New(ErrInvalidUserType)
	}

	if !isValidUserStatus(user.Status) {
		return User{}, errors.New(ErrInvalidUserStatus)
	}

	return user, nil
}

func isValidEmail(email string) bool {
	return utils.IsValidEmail(email)
}

func isValidUserType(userType string) bool {
	userTypes := constants.Values.USER_TYPES.Values()

	return utils.IsInSlice(userTypes, userType)
}

func isValidUserStatus(userStatus string) bool {
	userStatuses := constants.Values.USER_STATUSES.Values()

	return utils.IsInSlice(userStatuses, userStatus)
}
