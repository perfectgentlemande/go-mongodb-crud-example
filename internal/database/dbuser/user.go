package dbuser

import (
	"context"
	"fmt"
	"time"

	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/service"
	"go.mongodb.org/mongo-driver/bson"
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
func (d *Database) ListUsers(ctx context.Context, params *service.ListUsersParams) ([]service.User, error) {
	// TODO add options

	cur, err := d.userCollection.Find(ctx, bson.M{})
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
func (d *Database) UpdateUserByID(ctx context.Context, id string, user *service.User) error {
	// be careful: this way does not unset the old fields (this is actually closer to PATCH request implementation)
	// if you want full-document upload implementation, check docs for MongoDB
	// and probably you should get the existing document on the service layer first
	// to guarantee the full-document upload (keep data like created_at) without knowing how any DB does the update

	_, err := d.userCollection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"&set": userFromSevice(user)})

	if err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}

	return nil
}
func (d *Database) DeleteUserByID(ctx context.Context, id string) error {
	_, err := d.userCollection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	return nil
}
