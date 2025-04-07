package domain

import (
	"context"
	"time"
)

type Friendship struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	FriendID  string    `json:"friend_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FriendshipRepository interface {
	AddFriend(ctx context.Context, friendship *Friendship) (*Friendship, error)
	RemoveFriend(ctx context.Context, friendship *Friendship) error
	GetFriends(ctx context.Context, userID string) ([]*Friendship, error)
	GetByUserIDAndFriendID(ctx context.Context, userID string, friendID string) (*Friendship, error)
}
