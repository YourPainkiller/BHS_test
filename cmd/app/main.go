package main

import (
	"context"
	"log"

	"github.com/YourPainkiller/BHS_test/internal/repository"
	"github.com/YourPainkiller/BHS_test/internal/repository/postgres"
	"github.com/YourPainkiller/BHS_test/internal/serivice"
	"github.com/YourPainkiller/BHS_test/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
	serivice.Run(*storeUseCase)
}

func newStorage(pool *pgxpool.Pool) repository.Facade {
	txManager := postgres.NewTxManager(pool)
	pgRepo := postgres.NewPgRepository(txManager)
	return repository.NewStorageFacade(*pgRepo, txManager)
}
