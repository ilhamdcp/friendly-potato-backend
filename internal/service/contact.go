package service

import (
	"context"
	"errors"
	"time"

	"github.com/ilhamdcp/friendly-potato/internal/domain"
)

type ContactService interface {
	AddContact(ctx context.Context, friendship *domain.Contact) (*domain.Contact, error)
	RemoveContact(ctx context.Context, friendship *domain.Contact) error
	GetContacts(ctx context.Context, userID string) ([]*domain.Contact, error)
}

type contactServiceImpl struct {
	contactRepo domain.ContactRepository
}

func NewContactServiceImpl(contactRepo domain.ContactRepository) *contactServiceImpl {
	return &contactServiceImpl{
		contactRepo: contactRepo,
	}
}
func (fs *contactServiceImpl) AddContact(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
	if contact.Username == "" {
		return nil, errors.New("user id cannot be empty")
	}
	if contact.ContactUsername == "" {
		return nil, errors.New("friend id cannot be empty")
	}
	existingContact, err := fs.contactRepo.GetByUsernameAndContactUsername(ctx, contact.Username, contact.ContactUsername)
	if err != nil {
		return nil, err
	}
	if existingContact != nil {
		return existingContact, nil
	}
	contact.CreatedAt = time.Now()

	userContact, err := fs.contactRepo.AddContact(ctx, contact)
	if err != nil {
		return nil, err
	}

	userID := contact.Username
	contact.Username = contact.ContactUsername
	contact.ContactUsername = userID

	existingContact, err = fs.contactRepo.GetByUsernameAndContactUsername(ctx, contact.Username, contact.ContactUsername)
	if existingContact == nil {
		fs.contactRepo.AddContact(ctx, contact)
	}

	return userContact, nil

}
func (fs *contactServiceImpl) RemoveContact(ctx context.Context, friendship *domain.Contact) error {
	return fs.contactRepo.RemoveContact(ctx, friendship)
}
func (fs *contactServiceImpl) GetContacts(ctx context.Context, userID string) ([]*domain.Contact, error) {
	return fs.contactRepo.GetContacts(ctx, userID)
}
