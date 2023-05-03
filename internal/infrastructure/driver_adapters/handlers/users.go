package handlers

import (
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	"github.com/Edigiraldo/car-rent/pkg/httphandler"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Users struct {
	UsersService ports.UsersService
}

func NewUsers(us ports.UsersService) Users {
	return Users{
		UsersService: us,
	}
}

func (uh Users) SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser domain.User
	user, err := dtos.UserFromBody(r.Body)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if newUser, err = uh.UsersService.Register(r.Context(), user.ToDomain()); err != nil {
		if err.Error() == services.ErrEmailAlreadyRegistered {
			httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		log.Println(err)

		return
	}

	user.FromDomain(newUser)
	httphandler.WriteSuccessResponse(w, http.StatusCreated, user)
}

func (uh Users) Get(w http.ResponseWriter, r *http.Request) {
	var user dtos.User

	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	du, err := uh.UsersService.Get(r.Context(), ID)
	if err != nil {
		if err.Error() == services.ErrUserNotFound {
			httphandler.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
			log.Println(err)
		}

		return
	}

	user.FromDomain(du)
	httphandler.WriteSuccessResponse(w, http.StatusOK, user)
}

func (uh Users) FullUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	user, err := dtos.UserFromBody(r.Body)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get the ID from path param
	user.ID = ID

	if err = uh.UsersService.FullUpdate(r.Context(), user.ToDomain()); err != nil {
		if err.Error() == services.ErrUserNotFound {
			httphandler.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else if err.Error() == services.ErrEmailAlreadyRegistered {
			httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
			log.Println(err)
		}

		return
	}

	httphandler.WriteSuccessResponse(w, http.StatusOK, user)
}

func (uh Users) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	err = uh.UsersService.Delete(r.Context(), ID)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		log.Println(err)

		return
	}

	httphandler.WriteSuccessResponse(w, http.StatusNoContent, nil)
}
