package domain

import "context"

type User struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	UserPin           string `json:userPin`
	Token             string `json:"token"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUserName(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
}
