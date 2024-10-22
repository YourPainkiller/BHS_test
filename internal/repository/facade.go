package repository

import (
	"context"

	"github.com/YourPainkiller/BHS_test/internal/dto"
	"github.com/YourPainkiller/BHS_test/internal/repository/postgres"
)

type Facade interface {
	AddUser(ctx context.Context, req dto.UserDto) error
	GetUserByUsername(ctx context.Context, username string) (*dto.UserDto, error)
	AddRefreshSession(ctx context.Context, req dto.RefreshSessionDto) error
	AddAsset(ctx context.Context, req dto.AssetDto) error
	DeleteAsset(ctx context.Context, req dto.DeleteAssetDto) error
	GetAssetInfo(ctx context.Context, assetName string) (*dto.AssetDto, error)
	Refresh(ctx context.Context, req dto.UpdateRefreshDto) error
}

type storageFacade struct {
	txManager    postgres.TransactionManager
	pgRepository postgres.PgRepository
}

func NewStorageFacade(pgRepository postgres.PgRepository, txManager postgres.TransactionManager) *storageFacade {
	return &storageFacade{pgRepository: pgRepository, txManager: txManager}
}

func (s *storageFacade) AddUser(ctx context.Context, req dto.UserDto) error {
	return s.txManager.RunReadWriteCommited(ctx, func(ctxTx context.Context) error {
		err := s.pgRepository.AddUser(ctxTx, req)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *storageFacade) GetUserByUsername(ctx context.Context, username string) (*dto.UserDto, error) {
	user, err := s.pgRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *storageFacade) AddRefreshSession(ctx context.Context, req dto.RefreshSessionDto) error {
	err := s.pgRepository.AddRefreshSession(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *storageFacade) AddAsset(ctx context.Context, req dto.AssetDto) error {
	err := s.pgRepository.AddAsset(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *storageFacade) DeleteAsset(ctx context.Context, req dto.DeleteAssetDto) error {
	err := s.pgRepository.DeleteAsset(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *storageFacade) GetAssetInfo(ctx context.Context, assetName string) (*dto.AssetDto, error) {
	asset, err := s.pgRepository.GetAssetInfo(ctx, assetName)
	if err != nil {
		return nil, err
	}

	return asset, nil
}
func (s *storageFacade) Refresh(ctx context.Context, req dto.UpdateRefreshDto) error {
	err := s.pgRepository.Refresh(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
