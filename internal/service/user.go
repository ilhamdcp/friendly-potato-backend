package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ilhamdcp/friendly-potato/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userRepo   domain.UserRepository
	hashSecret string
}

func NewUserServiceImpl(userRepo domain.UserRepository, hashSecret string) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:   userRepo,
		hashSecret: hashSecret,
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "ilhamdcp",
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30 * 3)),
	})

	tokenString, err := token.SignedString(token)

	user.Token = tokenString

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

func (us *UserServiceImpl) SignInUser(ctx context.Context, user *domain.User) (map[string]string, error) {
	if user.Username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if user.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	existingUser, err := us.userRepo.GetByUserName(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": existingUser.Username,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	tokenString, err := token.SignedString([]byte(us.hashSecret))
	if err != nil {
		return nil, err
	}

	existingUser.Token = tokenString

	err = us.userRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}
	return map[string]string{"token": existingUser.Token}, nil
}

func (us *UserServiceImpl) SignOutUser(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, errors.New("username cannot be empty")
	}

	user, err := us.userRepo.GetByUserName(ctx, username)
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, errors.New("user not found")
	}

	user.Token = ""
	err = us.userRepo.Update(ctx, user)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (us *UserServiceImpl) AuthenticateUser(ctx context.Context, token string) bool {
	if token == "" {
		return false
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims := &jwt.RegisteredClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return false
	}

	if !tkn.Valid {
		return false
	}

	user, err := us.userRepo.GetByUserName(ctx, claims.Subject)
	if err != nil {
		return false
	}

	if user == nil {
		return false
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return false
	}
	if exp.Before(time.Now()) {
		user.Token = ""
		us.userRepo.Update(ctx, user)
		return false
	}

	return true
}
