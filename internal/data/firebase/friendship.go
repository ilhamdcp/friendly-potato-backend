package firebase

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	. "github.com/ilhamdcp/friendly-potato/internal/domain"
)

type FriendshipRepositoryImpl struct {
	client *firestore.Client
}

func NewFriendshipRepository(client *firestore.Client) *FriendshipRepositoryImpl {
	return &FriendshipRepositoryImpl{
		client: client,
	}

}

func (fr *FriendshipRepositoryImpl) AddFriend(ctx context.Context, friendship *Friendship) (*Friendship, error) {
	ref := fr.client.Collection("friendships").NewDoc()
	friendship.ID = ref.ID

	_, err := ref.Set(ctx, friendship)

	if err != nil {
		return nil, err
	}

	snapshot, _ := fr.client.Collection("users").Doc(ref.ID).Get(ctx)
	var newFriendship Friendship
	if err := snapshot.DataTo(&newFriendship); err != nil {
		return nil, err
	}
	return &newFriendship, nil
}

func (fr *FriendshipRepositoryImpl) RemoveFriend(ctx context.Context, friendship *Friendship) error {
	return errors.New("not implemented")
}

func (fr *FriendshipRepositoryImpl) GetFriends(ctx context.Context, userID string) ([]*Friendship, error) {
	snapshot, err := fr.client.Collection("users").Where("UserID", "==", userID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var friendships []*Friendship
	for _, doc := range snapshot {
		var friendship Friendship
		if err := doc.DataTo(&friendship); err != nil {
			return nil, err
		}
		friendships = append(friendships, &friendship)
	}
	return friendships, nil
}

func (fr *FriendshipRepositoryImpl) GetByUserIDAndFriendID(ctx context.Context, userID string, friendID string) (*Friendship, error) {
	snapshot, err := fr.client.Collection("friendships").Where("UserID", "==", userID).Where("FriendID", "==", friendID).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(snapshot) == 0 {
		return nil, nil
	}
	var friendship Friendship
	if err := snapshot[0].DataTo(&friendship); err != nil {
		return nil, err
	}
	return &friendship, nil
}
