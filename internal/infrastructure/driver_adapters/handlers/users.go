package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
)

type Users struct {
	UsersService ports.UsersService
}

func NewUsers(us ports.UsersService) *Users {
	return &Users{
		UsersService: us,
	}
}

func (us *Users) SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser domain.User
	user, err := dtos.UserFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if newUser, err = us.UsersService.Register(r.Context(), user.ToDomain()); err != nil {
		log.Println(err)
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)

		return
	}

	user.FromDomain(newUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
