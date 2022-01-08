package user

import (
	"context"
	"fmt"
	"project/pkg/logging"
)

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(userStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: userStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	Create(ctx context.Context, user User) (userID string, err error)
	GetUserFriends(ctx context.Context, userID string) (friends []string, err error)
	UpdateAge(ctx context.Context, id string, age string) error
	Delete(ctx context.Context, userID string) error
	MakeFriends(ctx context.Context, firstUserID string, secondUserID string) (firstUser User, secondUser User, err error)
}

func (s service) Create(ctx context.Context, user User) (userID string, err error) {
	s.logger.Info("create user")
	userID, err = s.storage.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user. error: %w", err)
	}
	return userID, nil
}

func (s service) GetUserFriends(ctx context.Context, userID string) (friends []string, err error) {
	friends, err = s.storage.GetUserFriends(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user. error: %w", err)
	}
	return friends, nil
}

func (s service) UpdateAge(ctx context.Context, id string, age string) error {
	err := s.storage.UpdateAge(ctx, id, age)
	if err != nil {
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return err
}

func (s service) Delete(ctx context.Context, userID string) error {
	err := s.storage.Delete(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return err
}

func (s service) MakeFriends(ctx context.Context, firstUserID string, secondUserID string) (firstUser User, secondUser User, err error) {
	firstUser, secondUser, err = s.storage.MakeFriends(ctx, firstUserID, secondUserID)
	if err != nil {
		return firstUser, secondUser, fmt.Errorf("failed to make friends. error: %w", err)
	}
	return firstUser, secondUser, nil
}
