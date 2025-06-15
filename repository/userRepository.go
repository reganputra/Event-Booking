package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-rest-api/model"
	"go-rest-api/utils"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Validate(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (s *userRepository) Create(ctx context.Context, u *model.User) error {
	query := "INSERT INTO users (email, password) VALUES (?, ?)"
	stmt, err := s.db.PrepareContext(ctx, query)
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

func (s *userRepository) Validate(ctx context.Context, u *model.User) error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := s.db.QueryRowContext(ctx, query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.Id, &retrievedPassword)
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
