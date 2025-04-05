package service

import (
	"context"
	"errors"

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

func (us *UserServiceImpl) CreateUser(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		return errors.New("User ID cannot be empty")
	}
	if user.Username == "" {
		return errors.New("Username cannot be empty")
	}
	return us.userRepo.Create(ctx, user)
}
func (us *UserServiceImpl) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("User ID cannot be empty")
	}
	return us.userRepo.GetByID(ctx, id)
}
func (us *UserServiceImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	return us.userRepo.Update(ctx, user)
}
func (us *UserServiceImpl) SignInWithGoogle(ctx context.Context, googleID string) (*domain.User, error) {
	if googleID == "" {
		return nil, errors.New("google ID is required")
	}

	existingUser, err := us.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return existingUser, nil // User already exists
	}

	// User doesn't exist, create a new user
	newUser := &domain.User{
		GoogleID: googleID,
		Username: "New User", //Provide a default username.
	}

	err = us.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
