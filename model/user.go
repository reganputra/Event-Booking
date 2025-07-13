package model

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `binding:"required,email"`
	Password string    `binding:"required,min=8"`
	Role     string    `binding:"omitempty,oneof=user admin"`
}
