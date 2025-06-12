package model

import (
	"context"
	"go-rest-api/connection"
	"go-rest-api/utils"
)

type User struct {
	Id       int64
	Email    string `bind:"required"`
	Password string `bind:"required"`
}

func (u *User) CreateUser(ctx context.Context) error {
	query := "INSERT INTO users (email, password) VALUES (?, ?)"
	stmt, err := connection.DB.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, u.Email, utils.HashPassword(u.Password))
	if err != nil {
		panic(err)
	}

	userId, err := result.LastInsertId()

	u.Id = userId
	return err
}
