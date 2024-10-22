package usecase

import (
	"context"
	"testing"

	"github.com/YourPainkiller/BHS_test/internal/domain"
	"github.com/YourPainkiller/BHS_test/internal/dto"
	"github.com/YourPainkiller/BHS_test/internal/usecase/mock"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ctrl := minimock.NewController(t)
	facadeMock := mock.NewFacadeMock(ctrl)
	uc := &StoreUseCase{psqlRepo: facadeMock}

	facadeMock.AddUserMock.Expect(context.TODO(), dto.UserDto{UserName: "admin", UserPassword: "admin"}).Return(nil)
	assert.NoError(t, uc.RegisterUser(context.TODO(), dto.UserDto{UserName: "admin", UserPassword: "admin"}))

	facadeMock.AddUserMock.Expect(context.TODO(), dto.UserDto{UserName: "admin", UserPassword: "admin"}).Return(domain.ErrAlreadyExists)
	assert.ErrorIs(t, uc.RegisterUser(context.TODO(), dto.UserDto{UserName: "admin", UserPassword: "admin"}), domain.ErrAlreadyExists)
}

func TestLogin(t *testing.T) {
	ctrl := minimock.NewController(t)
	facadeMock := mock.NewFacadeMock(ctrl)
	uc := &StoreUseCase{psqlRepo: facadeMock}

	//correct/wrong cred
	facadeMock.GetUserByUsernameMock.Expect(context.TODO(), "admin").Return(&dto.UserDto{UserName: "admin", UserPassword: "admin"}, nil)
	_, err := uc.LoginUser(context.TODO(), dto.UserDto{UserName: "admin", UserPassword: "admin"})
	assert.NoError(t, err)
	_, err = uc.LoginUser(context.TODO(), dto.UserDto{UserName: "admin", UserPassword: "user"})
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)

	//No such user
	facadeMock.GetUserByUsernameMock.Expect(context.TODO(), "admin").Return(&dto.UserDto{UserName: "admin", UserPassword: "admin"}, pgx.ErrNoRows)
	_, err = uc.LoginUser(context.TODO(), dto.UserDto{UserName: "admin", UserPassword: "admin"})
	assert.ErrorIs(t, err, domain.ErrNoSuchUser)
}

func TestAddAsset(t *testing.T) {
	ctrl := minimock.NewController(t)
	facadeMock := mock.NewFacadeMock(ctrl)
	uc := &StoreUseCase{psqlRepo: facadeMock}

	//successfull add
	facadeMock.GetAssetInfoMock.Expect(context.TODO(), "tree").Return(&dto.AssetDto{}, pgx.ErrNoRows)
	facadeMock.AddAssetMock.Expect(context.TODO(), dto.AssetDto{UserId: 1, AssetName: "tree", AssetDescr: "lorem ipsum", AssetPrice: 123}).Return(nil)
	uc.AddAsset(context.TODO(), dto.AssetDto{UserId: 1, AssetName: "tree", AssetDescr: "lorem ipsum", AssetPrice: 123})

	//already exists
	facadeMock.GetAssetInfoMock.Expect(context.TODO(), "tree").Return(&dto.AssetDto{}, nil)
	assert.ErrorIs(t, uc.AddAsset(context.TODO(), dto.AssetDto{UserId: 1, AssetName: "tree", AssetDescr: "lorem ipsum", AssetPrice: 123}), domain.ErrAlreadyExists)
}

func TestDeleteAsset(t *testing.T) {
	ctrl := minimock.NewController(t)
	facadeMock := mock.NewFacadeMock(ctrl)
	uc := &StoreUseCase{psqlRepo: facadeMock}

	//successfull delete
	facadeMock.DeleteAssetMock.Expect(context.TODO(), dto.DeleteAssetDto{AssetName: "tree"}).Return(nil)
	assert.NoError(t, uc.DeleteAsset(context.TODO(), dto.DeleteAssetDto{AssetName: "tree"}))

	//no such order
	facadeMock.DeleteAssetMock.Expect(context.TODO(), dto.DeleteAssetDto{AssetName: "tree"}).Return(domain.ErrNoSuchAsset)
	assert.ErrorIs(t, uc.DeleteAsset(context.TODO(), dto.DeleteAssetDto{AssetName: "tree"}), domain.ErrNoSuchAsset)
}

func TestBuyAsset(t *testing.T) {
	ctrl := minimock.NewController(t)
	facadeMock := mock.NewFacadeMock(ctrl)
	uc := &StoreUseCase{psqlRepo: facadeMock}

	//successfull buy
	facadeMock.GetAssetInfoMock.Expect(context.TODO(), "tree").Return(&dto.AssetDto{AssetName: "tree"}, nil)
	_, err := uc.BuyAsset(context.TODO(), dto.BuyAssetDto{AssetName: "tree"})
	assert.NoError(t, err)

	//No such order
	facadeMock.GetAssetInfoMock.Expect(context.TODO(), "tree").Return(&dto.AssetDto{AssetName: "tree"}, pgx.ErrNoRows)
	_, err = uc.BuyAsset(context.TODO(), dto.BuyAssetDto{AssetName: "tree"})
	assert.ErrorIs(t, err, domain.ErrNoSuchAsset)
}
