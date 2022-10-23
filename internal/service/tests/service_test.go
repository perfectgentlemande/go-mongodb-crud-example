package tests

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/service"
	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/service/tests/mock"
)

// these tests are made for the sake of checking how to use gomock
func getMockedUser(num int) service.User {
	return service.User{
		ID:        uuid.NewString(),
		Username:  "someUsername" + strconv.Itoa(num),
		Email:     "some_email" + strconv.Itoa(num) + "@some.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
func getMockedUsers(count int) []service.User {
	res := make([]service.User, 0, count)

	for i := 1; i <= count; i++ {
		res = append(res, getMockedUser(i))
	}

	return res
}

func TestListUsers(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserStorage := mock.NewMockUserStorage(mockCtrl)
	testSrvc := service.NewService(mockUserStorage)

	mockedUsers := getMockedUsers(50)
	mockUserStorage.EXPECT().ListUsers(ctx, &service.ListUsersParams{}).Return(mockedUsers, nil).Times(1)

	listedUsers, err := testSrvc.ListUsers(ctx, &service.ListUsersParams{})
	if err != nil {
		t.Error(err)
	}

	if len(listedUsers) != len(mockedUsers) {
		t.Errorf("expected listed %d users, got: %d", len(listedUsers), len(mockedUsers))
	}
}
func TestCreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
}
func TestGetUserByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
}
func TestUpdateUserByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
}
func TestDeleteUserByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
}
