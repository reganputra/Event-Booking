package model

type User struct {
	Id       int64
	Email    string `bind:"required"`
	Password string `bind:"required"`
}
