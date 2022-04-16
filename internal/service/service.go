package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Service struct {
	userStorage UserStorage
}

func NewService(userStorage UserStorage) *Service {
	return &Service{userStorage: userStorage}
}

type ListUsersParams struct {
	Limit  *int64
	Offset *int64
}

func (s *Service) ListUsers(ctx context.Context, params *ListUsersParams) ([]User, error) {
	return s.userStorage.ListUsers(ctx, params)
}
func (s *Service) CreateUser(ctx context.Context, user *User) (string, error) {
	now := time.Now()

	user.ID = uuid.NewString()
	user.CreatedAt, user.UpdatedAt = now, now

	return s.userStorage.CreateUser(ctx, user)
}
func (s *Service) GetUserByID(ctx context.Context, id string) (User, error) {
	return s.userStorage.GetUserByID(ctx, id)
}
func (s *Service) UpdateUserByID(ctx context.Context, id string, user *User) (User, error) {
	user.ID = id
	user.UpdatedAt = time.Now()

	return s.userStorage.UpdateUserByID(ctx, id, user)
}
func (s *Service) DeleteUserByID(ctx context.Context, id string) error {
	return s.userStorage.DeleteUserByID(ctx, id)
}
