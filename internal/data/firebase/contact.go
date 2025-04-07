package firebase

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	. "github.com/ilhamdcp/friendly-potato/internal/domain"
)

type ContactRepositoryImpl struct {
	client *firestore.Client
}

func NewContactRepository(client *firestore.Client) *ContactRepositoryImpl {
	return &ContactRepositoryImpl{
		client: client,
	}

}

func (fr *ContactRepositoryImpl) AddContact(ctx context.Context, contact *Contact) (*Contact, error) {
	ref := fr.client.Collection("contacts").NewDoc()
	contact.ID = ref.ID

	_, err := ref.Set(ctx, contact)

	if err != nil {
		return nil, err
	}

	snapshot, _ := fr.client.Collection("contacts").Doc(ref.ID).Get(ctx)
	var newContact Contact
	if err := snapshot.DataTo(&newContact); err != nil {
		return nil, err
	}
	return &newContact, nil
}

func (fr *ContactRepositoryImpl) RemoveContact(ctx context.Context, contact *Contact) error {
	return errors.New("not implemented")
}

func (fr *ContactRepositoryImpl) GetContacts(ctx context.Context, username string) ([]*Contact, error) {
	snapshot, err := fr.client.Collection("contacts").Where("Username", "==", username).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var contacts []*Contact
	for _, doc := range snapshot {
		var contact Contact
		if err := doc.DataTo(&contact); err != nil {
			return nil, err
		}
		contacts = append(contacts, &contact)
	}
	return contacts, nil
}

func (fr *ContactRepositoryImpl) GetByUsernameAndContactUsername(ctx context.Context, username string, contactUsername string) (*Contact, error) {
	snapshot, err := fr.client.Collection("contacts").Where("Username", "==", username).Where("ContactUsername", "==", contactUsername).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(snapshot) == 0 {
		return nil, nil
	}
	var contact Contact
	if err := snapshot[0].DataTo(&contact); err != nil {
		return nil, err
	}
	return &contact, nil
}
