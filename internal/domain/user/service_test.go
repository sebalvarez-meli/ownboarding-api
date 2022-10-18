package user

import (
	"context"
	"github.com/api_base/internal/domain"
	"github.com/api_base/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type fakeContainer struct {
	domain.Container
	UserRepoMock  *userRepositoryMock
	TokenRepoMock *tokenRepositoryMock
}

func newContainerMock() *fakeContainer {
	fc := &fakeContainer{
		UserRepoMock:  &userRepositoryMock{},
		TokenRepoMock: &tokenRepositoryMock{},
	}
	fc.Container = domain.Container{
		UserRepo:  fc.UserRepoMock,
		TokenRepo: fc.TokenRepoMock,
	}
	return fc
}

type userRepositoryMock struct {
	mock.Mock
}

func (mr *userRepositoryMock) Get(ctx context.Context, id string) (*model.User, error) {
	args := mr.Called(ctx, id)
	res := args.Get(0).(*model.User)
	return res, args.Error(1)
}

type tokenRepositoryMock struct {
	mock.Mock
}

func (tk *tokenRepositoryMock) Get(ctx context.Context, id string) (model.Token, error) {
	args := tk.Called(ctx, id)
	res := args.Get(0).(model.Token)
	return res, args.Error(1)
}

func initTest() (context.Context, *fakeContainer, Service) {
	ctn := newContainerMock()
	srv := NewService(ctn.Container)
	return context.Background(), ctn, srv
}

func TestService_Get(t *testing.T) {
	ctx, cnt, srv := initTest()

	userDb := &model.User{Id: "1"}
	tokenResp := model.Token{Id: "token_1", UserId: "1"}
	cnt.UserRepoMock.On("Get", ctx, "1").Return(userDb, nil)
	cnt.TokenRepoMock.On("Get", ctx, "1").Return(tokenResp, nil)

	user, err := srv.Get(ctx, "1")

	assert.Nil(t, err)
	assert.Equal(t, "1", user.Id)
	assert.Equal(t, "token_1", user.Token.Id)
}
