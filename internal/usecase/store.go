package usecase

import (
	"context"

	"github.com/YourPainkiller/BHS_test/internal/domain"
	"github.com/YourPainkiller/BHS_test/internal/dto"
	"github.com/YourPainkiller/BHS_test/internal/repository"
)

type StoreUseCase struct {
	psqlRepo repository.Facade
}

func NewStoreUseCase(psqlRepo repository.Facade) *StoreUseCase {
	return &StoreUseCase{psqlRepo: psqlRepo}
}

func (sc *StoreUseCase) RegisterUser(ctx context.Context, req dto.RegisterUserDto) error {
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
