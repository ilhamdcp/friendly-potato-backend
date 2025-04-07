package domain

import (
	"context"
	"time"
)

type Contact struct {
	ID              string    `json:"id"`
	Username        string    `json:"username"`
	ContactUsername string    `json:"contactUsername"`
	CreatedAt       time.Time `json:"createdAt"`
}

type ContactRepository interface {
	AddContact(ctx context.Context, contact *Contact) (*Contact, error)
	RemoveContact(ctx context.Context, contact *Contact) error
	GetContacts(ctx context.Context, username string) ([]*Contact, error)
	GetByUsernameAndContactUsername(ctx context.Context, username string, contactUsername string) (*Contact, error)
}
