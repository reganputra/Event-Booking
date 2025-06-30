package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-rest-api/model"
	"go-rest-api/utils"
	"strings"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Validate(ctx context.Context, user *model.User) error
	GetAll(ctx context.Context) ([]model.User, error)
	GetById(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (s *userRepository) Create(ctx context.Context, u *model.User) error {
	if u.Role == "" {
		u.Role = "user"
	}

	query := "INSERT INTO users (email, password) VALUES ($1, $2)"
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, u.Email, utils.HashPassword(u.Password))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("email already registered")
		}
		return err // Don't panic, just return the error
	}

	userId, err := result.LastInsertId()

	u.Id = userId
	return err
}

func (s *userRepository) GetAll(ctx context.Context) ([]model.User, error) {
	query := "SELECT id, email, role FROM users"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Email, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userRepository) GetById(ctx context.Context, id int64) (*model.User, error) {
	query := "SELECT id, email, role FROM users WHERE id = $1"
	row := s.db.QueryRowContext(ctx, query, id)

	var user model.User
	err := row.Scan(&user.Id, &user.Email, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := "SELECT id, email, role FROM users WHERE email = $1"
	row := s.db.QueryRowContext(ctx, query, email)

	var user model.User
	err := row.Scan(&user.Id, &user.Email, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *userRepository) Update(ctx context.Context, u *model.User) error {
	var updateQuery string
	var args []interface{}

	if u.Password != "" {
		updateQuery = "UPDATE users SET email = $1, password = $2, role = $4 WHERE id = $5"
		args = []interface{}{u.Email, utils.HashPassword(u.Password), u.Role, u.Id}
	} else {
		updateQuery = "UPDATE users SET email = $1, role = $2 WHERE id = $3"
		args = []interface{}{u.Email, u.Role, u.Id}
	}

	stmt, err := s.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return errors.New("email already registered")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found or no changes made")
	}
	return nil
}

func (s *userRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (s *userRepository) Validate(ctx context.Context, u *model.User) error {
	query := "SELECT id, password, role FROM users WHERE email = $1"
	row := s.db.QueryRowContext(ctx, query, u.Email)

	var retrievedPassword string
	var retrievedRole string
	err := row.Scan(&u.Id, &retrievedPassword, &retrievedRole)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	if !passwordIsValid {
		return errors.New("invalid Credentials")
	}
	u.Password = retrievedPassword
	u.Role = retrievedRole
	return nil
}
