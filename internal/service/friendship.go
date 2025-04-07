package service

import (
	"context"
	"errors"
	"time"

	"github.com/ilhamdcp/friendly-potato/internal/domain"
)

type FriendshipService interface {
	AddFriend(ctx context.Context, friendship *domain.Friendship) (*domain.Friendship, error)
	RemoveFriend(ctx context.Context, friendship *domain.Friendship) error
	GetFriends(ctx context.Context, userID string) ([]*domain.Friendship, error)
}

type friendshipServiceImpl struct {
	friendshipRepo domain.FriendshipRepository
}

func NewFriendshipServiceImpl(friendshipRepo domain.FriendshipRepository) *friendshipServiceImpl {
	return &friendshipServiceImpl{
		friendshipRepo: friendshipRepo,
	}
}
func (fs *friendshipServiceImpl) AddFriend(ctx context.Context, friendship *domain.Friendship) (*domain.Friendship, error) {
	if friendship.UserID == "" {
		return nil, errors.New("user id cannot be empty")
	}
	if friendship.FriendID == "" {
		return nil, errors.New("friend id cannot be empty")
	}
	existingFriendship, err := fs.friendshipRepo.GetByUserIDAndFriendID(ctx, friendship.UserID, friendship.FriendID)
	if err != nil {
		return nil, err
	}
	if existingFriendship != nil {
		return existingFriendship, nil
	}
	friendship.CreatedAt = time.Now()

	userFriendship, err := fs.friendshipRepo.AddFriend(ctx, friendship)
	if err != nil {
		return nil, err
	}

	userID := friendship.UserID
	friendship.UserID = friendship.FriendID
	friendship.FriendID = userID

	existingFriendship, err = fs.friendshipRepo.GetByUserIDAndFriendID(ctx, friendship.UserID, friendship.FriendID)
	if existingFriendship == nil {
		fs.friendshipRepo.AddFriend(ctx, friendship)
	}

	return userFriendship, nil

}
func (fs *friendshipServiceImpl) RemoveFriend(ctx context.Context, friendship *domain.Friendship) error {
	return fs.friendshipRepo.RemoveFriend(ctx, friendship)
}
func (fs *friendshipServiceImpl) GetFriends(ctx context.Context, userID string) ([]*domain.Friendship, error) {
	return fs.friendshipRepo.GetFriends(ctx, userID)
}
