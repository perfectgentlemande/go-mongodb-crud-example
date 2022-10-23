package service

import "context"

//go:generate mockgen -source=storage.go -destination=./tests/mock/storage_mock.go -package=mock

type UserStorage interface {
	ListUsers(ctx context.Context, params *ListUsersParams) ([]User, error)
	CreateUser(ctx context.Context, user *User) (string, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	UpdateUserByID(ctx context.Context, id string, user *User) (User, error)
	DeleteUserByID(ctx context.Context, id string) error
}
