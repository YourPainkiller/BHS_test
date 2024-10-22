package usecase

import (
	"context"
	"errors"

	"github.com/YourPainkiller/BHS_test/internal/domain"
	"github.com/YourPainkiller/BHS_test/internal/dto"
	"github.com/YourPainkiller/BHS_test/internal/repository"
	"github.com/YourPainkiller/BHS_test/internal/repository/postgres"
	"github.com/jackc/pgx/v5"
)

type StoreUseCase struct {
	psqlRepo repository.Facade
}

func NewStoreUseCase(psqlRepo repository.Facade) *StoreUseCase {
	return &StoreUseCase{psqlRepo: psqlRepo}
}

func (sc *StoreUseCase) RegisterUser(ctx context.Context, req dto.UserDto) error {
	newUser, err := domain.NewUser(req.UserName, req.UserPassword)
	if err != nil {
		return err
	}

	err = sc.psqlRepo.AddUser(ctx, newUser.ToDTO())
	if postgres.UnwrapPgCode(err) == domain.UniqueErrCode {
		return domain.ErrAlreadyExists
	}

	if err != nil {
		return err
	}
	return nil
}

func (sc *StoreUseCase) LoginUser(ctx context.Context, req dto.UserDto) (int, error) {
	user, err := sc.psqlRepo.GetUserByUsername(ctx, req.UserName)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, domain.ErrNoSuchUser
	}
	if err != nil {
		return 0, err
	}
	if user.UserPassword != req.UserPassword {
		return 0, domain.ErrInvalidCredentials
	}

	return user.UserId, nil
}

func (sc *StoreUseCase) SetSession(ctx context.Context, req dto.RefreshSessionDto) error {
	err := sc.psqlRepo.AddRefreshSession(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (sc *StoreUseCase) AddAsset(ctx context.Context, req dto.AssetDto) error {
	_, err := sc.psqlRepo.GetAssetInfo(ctx, req.AssetName)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		err = sc.psqlRepo.AddAsset(ctx, req)
		if err != nil {
			return err
		}
		return nil
	}
	return domain.ErrAlreadyExists

}

func (sc *StoreUseCase) DeleteAsset(ctx context.Context, req dto.DeleteAssetDto) error {
	err := sc.psqlRepo.DeleteAsset(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (sc *StoreUseCase) BuyAsset(ctx context.Context, req dto.BuyAssetDto) (*dto.AssetDto, error) {
	asset, err := sc.psqlRepo.GetAssetInfo(ctx, req.AssetName)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, domain.ErrNoSuchAsset
		default:
			return nil, err
		}
	}

	return asset, nil
}

func (sc *StoreUseCase) Refresh(ctx context.Context, req dto.UpdateRefreshDto) error {
	err := sc.psqlRepo.Refresh(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
