package user

// file for data storage interface

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	GetUserFriends(ctx context.Context, userID string) (friends []string, err error)
	UpdateAge(ctx context.Context, id string, age string) error
	Delete(ctx context.Context, id string) error
	MakeFriends(ctx context.Context, firstUserID string, secondUserID string) (firstUser User, secondUser User, err error)
}
