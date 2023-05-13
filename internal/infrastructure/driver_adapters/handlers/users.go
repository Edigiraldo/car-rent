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

// @Summary Register a new user
// @Description Register a new user with the provided information
// @ID register-user
// @Accept json
// @Produce json
// @Param user body docs.UserRequest true "User information (allowed types: Customer, Admin; allowed statuses: Active, Inactive)"
// @Success 201 {object} docs.UserResponse "Created user"
// @Failure 400 {object} docs.ErrorResponseBadRequest "Bad Request"
// @Failure 500 {object} docs.ErrorResponseInternalServer "Internal Server Error"
// @Tags Users
// @Router /users [post]
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

// @Summary Get a user
// @Description Get a user by UUID
// @ID get-user
// @Produce json
// @Param id path string true "User UUID" format(uuid)
// @Success 201 {object} docs.UserResponse "Obtained user"
// @Failure 400 {object} docs.ErrorResponseBadRequest "Bad Request"
// @Failure 404 {object} docs.ErrorResponseNotFound "Not Found"
// @Failure 500 {object} docs.ErrorResponseInternalServer "Internal Server Error"
// @Tags Users
// @Router /users/{id} [get]
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

// @Summary Update a user
// @Description Update a user by UUID
// @ID update-user
// @Accept json
// @Produce json
// @Param id path string true "User UUID" format(uuid)
// @Param user body docs.UserRequest true "User information (allowed types: Customer, Admin; allowed statuses: Active, Inactive)"
// @Success 201 {object} docs.UserResponse "Updated user"
// @Failure 400 {object} docs.ErrorResponseBadRequest "Bad Request"
// @Failure 404 {object} docs.ErrorResponseNotFound "Not Found"
// @Failure 500 {object} docs.ErrorResponseInternalServer "Internal Server Error"
// @Tags Users
// @Router /users/{id} [put]
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

// @Summary Delete a user
// @Description Delete a user by UUID
// @ID delete-user
// @Produce json
// @Param id path string true "User UUID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} docs.ErrorResponseBadRequest "Bad Request"
// @Failure 404 {object} docs.ErrorResponseNotFound "Not Found"
// @Failure 500 {object} docs.ErrorResponseInternalServer "Internal Server Error"
// @Tags Users
// @Router /users/{id} [delete]
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
		if err.Error() == services.ErrUserNotFound {
			httphandler.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
			log.Println(err)
		}

		return
	}

	httphandler.WriteSuccessResponse(w, http.StatusNoContent, nil)
}
