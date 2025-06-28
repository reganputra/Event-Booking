package services

import (
	"context"
	"errors"
	"go-rest-api/model"
	"go-rest-api/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	ValidateUser(ctx context.Context, user *model.User) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int64) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (e *userService) CreateUser(ctx context.Context, user *model.User) error {
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}

	existingUser, _ := e.userRepository.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	return e.userRepository.Create(ctx, user)
}

func (e *userService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return e.userRepository.GetAll(ctx)
}

func (e *userService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := e.userRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (e *userService) UpdateUser(ctx context.Context, user *model.User) error {
	if user.Role != "user" && user.Role != "admin" {
		return errors.New("invalid role")
	}

	existingUser, err := e.userRepository.GetById(ctx, user.Id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}
	return e.userRepository.Update(ctx, user)
}

func (e *userService) DeleteUser(ctx context.Context, id int64) error {
	return e.userRepository.Delete(ctx, id)
}

func (e *userService) ValidateUser(ctx context.Context, user *model.User) error {
	return e.userRepository.Validate(ctx, user)
}
