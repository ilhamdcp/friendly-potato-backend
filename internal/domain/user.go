package domain

import "context"

type User struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
	GoogleID          string `json:"googleId"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*User, error)
	Update(ctx context.Context, user *User) error
}
