package model

import (
	"context"
	"errors"
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

func (u *User) ValidateUser(ctx context.Context) error {
	query := "SELECT password FROM users WHERE email = ?"
	row := connection.DB.QueryRowContext(ctx, query, u.Email)

	var retrievedPassword string
	err := row.Scan(&retrievedPassword)
	if err != nil {
		return errors.New("user not found")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	if !passwordIsValid {
		return errors.New("invalid Credentials")
	}
	u.Password = retrievedPassword
	return nil
}
