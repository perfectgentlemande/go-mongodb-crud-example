package dbuser

import (
	"context"
	"fmt"
	"time"

	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        string    `bson:"_id"`
	Username  string    `bson:"username"`
	Email     string    `bson:"email,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (u *User) ToSevice() service.User {
	return service.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func userFromSevice(u *service.User) User {
	return User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (d *Database) CreateUser(ctx context.Context, user *service.User) (string, error) {
	_, err := d.userCollection.InsertOne(ctx, userFromSevice(user))

	if err != nil {
		return "", fmt.Errorf("cannot insert user: %w", err)
	}

	return user.ID, nil
}

func withListUsersParams(opts *options.FindOptions, params *service.ListUsersParams) *options.FindOptions {
	if params.Limit != nil {
		opts = opts.SetLimit(*params.Limit)
	}
	if params.Offset != nil {
		opts = opts.SetSkip(*params.Offset)
	}

	return opts
}

func (d *Database) ListUsers(ctx context.Context, params *service.ListUsersParams) ([]service.User, error) {
	opts := options.Find()
	opts = withListUsersParams(opts, params)

	cur, err := d.userCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot find users: %w", err)
	}
	defer cur.Close(ctx)

	res := make([]service.User, 0)
	for cur.Next(ctx) {
		usr := User{}

		err = cur.Decode(&usr)
		if err != nil {
			return nil, fmt.Errorf("cannot decode user: %w", err)
		}

		res = append(res, usr.ToSevice())
	}

	return res, nil
}
func (d *Database) GetUserByID(ctx context.Context, id string) (service.User, error) {
	usr := User{}

	err := d.userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&usr)
	if err != nil {
		return service.User{}, fmt.Errorf("cannot find user: %w", err)
	}

	return usr.ToSevice(), nil
}
func (d *Database) UpdateUserByID(ctx context.Context, id string, user *service.User) (service.User, error) {
	usr := userFromSevice(user)

	_, err := d.userCollection.UpdateOne(ctx, bson.M{"_id": id}, usr)
	if err != nil {
		return service.User{}, fmt.Errorf("cannot update user: %w", err)
	}

	return *user, nil
}
func (d *Database) DeleteUserByID(ctx context.Context, id string) error {
	_, err := d.userCollection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	return nil
}
