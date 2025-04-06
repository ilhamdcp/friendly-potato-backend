package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
	. "github.com/ilhamdcp/friendly-potato/internal/domain"
)

type UserRepositoryImpl struct {
	client *firestore.Client
}

func NewUserRepository(client *firestore.Client) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		client: client,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *User) (*User, error) {
	ref := r.client.Collection("users").NewDoc()
	user.ID = ref.ID
	_, err := ref.Set(ctx, user)
	if err != nil {
		return nil, err
	}

	snapshot, _ := r.client.Collection("users").Doc(ref.ID).Get(ctx)
	var newUser User
	if err := snapshot.DataTo(&newUser); err != nil {
		return nil, err
	}
	return &newUser, nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*User, error) {
	doc, err := r.client.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	user.Password = ""
	user.UserPin = ""
	return &user, nil
}

func (r *UserRepositoryImpl) GetByGoogleID(ctx context.Context, googleID string) (*User, error) {
	doc, err := r.client.Collection("users").Where("google_id", "==", googleID).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(doc) == 0 {
		return nil, nil
	}
	var user User
	if err := doc[0].DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *User) error {
	_, err := r.client.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) GetByUserName(ctx context.Context, username string) (*User, error) {
	doc, err := r.client.Collection("users").Where("Username", "==", username).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(doc) == 0 {
		return nil, nil
	}
	var user User
	if err := doc[0].DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
