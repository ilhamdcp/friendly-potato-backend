package domain

import "context"

type User struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	UserPin           string `json:userPin`
	ProfilePictureUrl string `json:"profilePictureUrl"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*User, error)
	Update(ctx context.Context, user *User) error
}
