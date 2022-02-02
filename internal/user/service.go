package user

// file for user-service functions

import (
	"context"
	"fmt"
	"project/pkg/logging"
)

// service struct with logging
type service struct {
	storage Storage
	logger  logging.Logger
}

// NewService - func for initialization user-service
func NewService(userStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: userStorage,
		logger:  logger,
	}, nil
}

// Service - interface for description user-service functions
type Service interface {
	Create(ctx context.Context, user User) (userID string, err error)
	GetUserFriends(ctx context.Context, userID string) (friends []string, err error)
	UpdateAge(ctx context.Context, id string, age string) error
	Delete(ctx context.Context, userID string) error
	MakeFriends(ctx context.Context, firstUserID string, secondUserID string) (firstUser User, secondUser User, err error)
}

// Create - func for creating user
func (s service) Create(ctx context.Context, user User) (userID string, err error) {
	s.logger.Info("create user")
	userID, err = s.storage.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user. error: %w", err)
	}
	return userID, nil
}

// GetUserFriends - func for get all friends from one user
func (s service) GetUserFriends(ctx context.Context, userID string) (friends []string, err error) {
	friends, err = s.storage.GetUserFriends(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user. error: %w", err)
	}
	return friends, nil
}

// UpdateAge - func for updating age of one user
func (s service) UpdateAge(ctx context.Context, id string, age string) error {
	err := s.storage.UpdateAge(ctx, id, age)
	if err != nil {
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return err
}

// Delete - func for deleting one user from database
func (s service) Delete(ctx context.Context, userID string) error {
	err := s.storage.Delete(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return err
}

// MakeFriends - func that create friendship between tow user by appending users names in friends array of each user
func (s service) MakeFriends(ctx context.Context, firstUserID string, secondUserID string) (firstUser User, secondUser User, err error) {
	firstUser, secondUser, err = s.storage.MakeFriends(ctx, firstUserID, secondUserID)
	if err != nil {
		return firstUser, secondUser, fmt.Errorf("failed to make friends. error: %w", err)
	}
	return firstUser, secondUser, nil
}
