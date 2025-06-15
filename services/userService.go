package services

import (
	"context"
	"go-rest-api/model"
	"go-rest-api/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	ValidateUser(ctx context.Context, user *model.User) error
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
	return e.userRepository.Create(ctx, user)
}
func (e *userService) ValidateUser(ctx context.Context, user *model.User) error {
	return e.userRepository.Validate(ctx, user)
}
