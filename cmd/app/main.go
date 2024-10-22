package main

import (
	"context"
	"log"

	"github.com/YourPainkiller/BHS_test/internal/repository"
	"github.com/YourPainkiller/BHS_test/internal/repository/postgres"
	"github.com/YourPainkiller/BHS_test/internal/service"
	"github.com/YourPainkiller/BHS_test/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

// @title Asset Store API
// @version 1.0
// @license.name xd
// @host localhost:4000
// @accept json
// @produce json
// @schemes http

func main() {
	const psqlDSN = "postgres://postgres:qwe@localhost:5434/store?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	storageFacade := newStorage(pool)

	storeUseCase := usecase.NewStoreUseCase(storageFacade)
	service.Run(*storeUseCase)
}

func newStorage(pool *pgxpool.Pool) repository.Facade {
	txManager := postgres.NewTxManager(pool)
	pgRepo := postgres.NewPgRepository(txManager)
	return repository.NewStorageFacade(*pgRepo, txManager)
}
