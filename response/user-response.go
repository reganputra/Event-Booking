package response

import "github.com/google/uuid"

type UserResponse struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}
