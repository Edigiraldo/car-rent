package docs

import "github.com/google/uuid"

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type UserRequest struct {
	FirstName string `json:"first_name" example:"Isaac"`
	LastName  string `json:"last_name" example:"Newton"`
	Email     string `json:"email" example:"isaac.newton@cam.ac.uk"`
	Type      string `json:"type" example:"Customer"`
	Status    string `json:"status" example:"Active"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty" example:"b6dcf3b3-ec0a-9f31-4379-4b8e7b94a387"`
	FirstName string    `json:"first_name" example:"Isaac"`
	LastName  string    `json:"last_name" example:"Newton"`
	Email     string    `json:"email" example:"isaac.newton@cam.ac.uk"`
	Type      string    `json:"type" example:"Customer"`
	Status    string    `json:"status" example:"Active"`
}
