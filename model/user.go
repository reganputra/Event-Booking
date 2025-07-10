package model

type User struct {
	Id       int64
	Email    string `binding:"required,email"`
	Password string `binding:"required,min=8"`
	Role     string `binding:"omitempty,oneof=user admin"`
}
