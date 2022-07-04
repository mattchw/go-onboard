package test

import (
	"context"
	"testing"

	"github.com/mattchw/go-onboard/models"
	"github.com/mattchw/go-onboard/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImplMock struct {
	mock.Mock
}

func (us *UserServiceImplMock) Create(ctx context.Context, payload models.User) *mongo.InsertOneResult {
	args := us.Called(ctx, payload)
	return args.Get(0).(*mongo.InsertOneResult)
}

func TestNewUserServiceImpl(t *testing.T) {
	// setup
	coll := &mongo.Collection{}
	// execute
	us := services.NewUserServiceImpl(coll)
	// verify
	require.NotNil(t, us)
}

func TestCreateUser(t *testing.T) {
	mock := new(UserServiceImplMock)

	ctx := context.Background()
	user := models.User{
		FirstName: "Matt",
		LastName:  "Chw",
		Bio:       "",
		Age:       0,
		Gender:    "Male",
	}
	mock.On("Create", ctx, user).Return(&mongo.InsertOneResult{})

	result := mock.Create(ctx, user)

	print(result)
	mock.AssertExpectations(t)
}
