package service

import "context"

type UserStorage interface {
	ListUsers(ctx context.Context, params *ListUsersParams) ([]User, error)
	CreateUser(ctx context.Context, user *User) (string, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	UpdateUserByID(ctx context.Context, id string, user *User) (User, error)
	DeleteUserByID(ctx context.Context, id string) error
}
