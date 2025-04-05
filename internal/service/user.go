package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/ilhamdcp/friendly-potato/internal/domain"
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

	hash := sha256.New()

	hash.Write([]byte(user.Password))
	user.Password = fmt.Sprintf("%x", hash.Sum(nil))

	if user.UserPin != "" {
		hash.Write([]byte(user.UserPin))
		user.UserPin = fmt.Sprintf("%x", hash.Sum(nil))
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
