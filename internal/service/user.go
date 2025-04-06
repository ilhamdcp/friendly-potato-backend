package service

import (
	"context"
	"errors"

	"github.com/ilhamdcp/friendly-potato/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userRepo domain.UserRepository
}

func NewUserServiceImpl(userRepo domain.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (us *UserServiceImpl) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if user.Username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if user.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	if user.UserPin != "" {
		hashedPin, err := bcrypt.GenerateFromPassword([]byte(user.UserPin), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.UserPin = string(hashedPin)
	}

	existingUser, _ := us.userRepo.GetByID(ctx, user.UserPin)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	result, err := us.userRepo.Create(ctx, user)
	return result, err
}

func (us *UserServiceImpl) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	return us.userRepo.GetByID(ctx, id)
}

func (us *UserServiceImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	return us.userRepo.Update(ctx, user)
}
