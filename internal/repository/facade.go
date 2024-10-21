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
	// GetOrderById(ctx context.Context, id int) (dto.OrderDto, error)
	// UpdateOrderInfo(ctx context.Context, req dto.OrderDto) error
	// GetOrdersByUserId(ctx context.Context, userId int) (dto.UserOrdersResponse, error)
	// GetUserReturns(ctx context.Context) (dto.UserReturnsResponse, error)
	// DropTable(ctx context.Context) error
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

// func (s *storageFacade) GetOrderById(ctx context.Context, id int) (dto.OrderDto, error) {
// 	order, err := s.pgRepository.GetOrderById(ctx, id)
// 	if err != nil {
// 		return dto.OrderDto{}, err
// 	}
// 	return order, nil
// }

// func (s *storageFacade) UpdateOrderInfo(ctx context.Context, req dto.OrderDto) error {
// 	err := s.pgRepository.UpdateOrderInfo(ctx, req)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *storageFacade) GetOrdersByUserId(ctx context.Context, userId int) (dto.UserOrdersResponse, error) {
// 	orders, err := s.pgRepository.GetOrdersByUserId(ctx, userId)
// 	if err != nil {
// 		return dto.UserOrdersResponse{}, err
// 	}

// 	return orders, nil
// }

// func (s *storageFacade) GetUserReturns(ctx context.Context) (dto.UserReturnsResponse, error) {
// 	orders, err := s.pgRepository.GetUserReturns(ctx)
// 	if err != nil {
// 		return dto.UserReturnsResponse{}, err
// 	}

// 	return orders, nil
// }

// func (s *storageFacade) DropTable(ctx context.Context) error {
// 	err := s.pgRepository.DropTable(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
