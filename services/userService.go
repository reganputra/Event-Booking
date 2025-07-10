package services

import (
	"context"
	"go-rest-api/apperrors"
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
		return apperrors.ErrInvalidInput
	}
	if user.Password == "" {
		return apperrors.ErrInvalidInput
	}

	// The check for existing user is now handled by the repository,
	// which returns apperrors.ErrAlreadyExists.
	return e.userRepository.Create(ctx, user)
}

func (e *userService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return e.userRepository.GetAll(ctx)
}

func (e *userService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	// The repository now correctly returns apperrors.ErrNotFound.
	return e.userRepository.GetById(ctx, id)
}

func (e *userService) UpdateUser(ctx context.Context, user *model.User) error {
	if user.Role != "user" && user.Role != "admin" {
		return apperrors.ErrInvalidInput
	}

	// The check for an existing user is handled by the repository,
	// which returns apperrors.ErrNotFound.
	return e.userRepository.Update(ctx, user)
}

func (e *userService) DeleteUser(ctx context.Context, id int64) error {
	return e.userRepository.Delete(ctx, id)
}

func (e *userService) ValidateUser(ctx context.Context, user *model.User) error {
	return e.userRepository.Validate(ctx, user)
}
