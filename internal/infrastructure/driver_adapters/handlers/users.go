package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Users struct {
	UsersService ports.UsersService
}

func NewUsers(us ports.UsersService) *Users {
	return &Users{
		UsersService: us,
	}
}

func (uh *Users) SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser domain.User
	user, err := dtos.UserFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if newUser, err = uh.UsersService.Register(r.Context(), user.ToDomain()); err != nil {
		log.Println(err)
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)

		return
	}

	user.FromDomain(newUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (uh *Users) Get(w http.ResponseWriter, r *http.Request) {
	var user dtos.User

	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, ErrInvalidID, http.StatusBadRequest)

		return
	}

	du, err := uh.UsersService.Get(r.Context(), ID)
	if err != nil {
		if err.Error() == services.ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
			log.Println(err)
		}

		return
	}

	user.FromDomain(du)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (uh *Users) FullUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, ErrInvalidID, http.StatusBadRequest)

		return
	}

	user, err := dtos.UserFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	// Get the ID from path param
	user.ID = ID

	if err = uh.UsersService.FullUpdate(r.Context(), user.ToDomain()); err != nil {
		if err.Error() == services.ErrUserNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
			log.Println(err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
