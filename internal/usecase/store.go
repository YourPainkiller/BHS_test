package usecase

import (
	"context"
	"errors"

	"github.com/YourPainkiller/BHS_test/internal/domain"
	"github.com/YourPainkiller/BHS_test/internal/dto"
	"github.com/YourPainkiller/BHS_test/internal/repository"
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
	if err != nil {
		return err
	}
	return nil
}

func (sc *StoreUseCase) LoginUser(ctx context.Context, req dto.UserDto) (int, error) {
	user, err := sc.psqlRepo.GetUserByUsername(ctx, req.UserName)
	if user.UserPassword != req.UserPassword || errors.Is(err, pgx.ErrNoRows) {
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
