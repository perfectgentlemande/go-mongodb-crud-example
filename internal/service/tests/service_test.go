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
	"github.com/stretchr/testify/assert"
)

// these tests are made for the sake of checking how to use gomock: there's just CRUD and no any logic almost at all
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
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserStorage := mock.NewMockUserStorage(mockCtrl)
	testSrvc := service.NewService(mockUserStorage)

	mockedUser := getMockedUser(1)
	mockUserStorage.EXPECT().CreateUser(ctx, &mockedUser).Return(mockedUser.ID, nil).Times(1)

	createdID, err := testSrvc.CreateUser(ctx, &mockedUser)
	if err != nil {
		t.Error(err)
	}

	if createdID == "" {
		t.Error("empty created ID")
	}
}
func TestGetUserByID(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserStorage := mock.NewMockUserStorage(mockCtrl)
	testSrvc := service.NewService(mockUserStorage)

	mockedUser := getMockedUser(1)
	mockUserStorage.EXPECT().GetUserByID(ctx, mockedUser.ID).Return(mockedUser, nil).Times(1)

	gotUser, err := testSrvc.GetUserByID(ctx, mockedUser.ID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, mockedUser, gotUser)
}
func TestUpdateUserByID(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserStorage := mock.NewMockUserStorage(mockCtrl)
	testSrvc := service.NewService(mockUserStorage)

	now := time.Now()
	mockedUser := getMockedUser(1)
	mockedUser.Username = "someUsername" + strconv.Itoa(2)
	mockedUser.Email = "some_email" + strconv.Itoa(2) + "@some.com"
	mockedUser.UpdatedAt = now
	mockUserStorage.EXPECT().UpdateUserByID(ctx, mockedUser.ID, &mockedUser).Return(mockedUser, nil).Times(1)

	gotUser, err := testSrvc.UpdateUserByID(ctx, mockedUser.ID, &mockedUser)
	if err != nil {
		t.Error(err)
	}

	if mockedUser.ID != gotUser.ID {
		t.Errorf("changed ID from: %s to: %s", mockedUser.ID, gotUser.ID)
	}
	if mockedUser.Username != gotUser.Username {
		t.Errorf("expected username: %s to: %s", mockedUser.Username, gotUser.Username)
	}
	if mockedUser.Email != gotUser.Email {
		t.Errorf("expected email: %s to: %s", mockedUser.Email, gotUser.Email)
	}
}
func TestDeleteUserByID(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserStorage := mock.NewMockUserStorage(mockCtrl)
	testSrvc := service.NewService(mockUserStorage)

	mockedUser := getMockedUser(1)
	mockUserStorage.EXPECT().DeleteUserByID(ctx, mockedUser.ID).Return(nil).Times(1)

	err := testSrvc.DeleteUserByID(ctx, mockedUser.ID)
	if err != nil {
		t.Error(err)
	}
}
