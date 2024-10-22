package service

import (
	"log"
	"net/http"

	"github.com/YourPainkiller/BHS_test/internal/routes"
	"github.com/YourPainkiller/BHS_test/internal/usecase"
)

func Run(storeUsecase usecase.StoreUseCase) {
	srv := &http.Server{
		Addr:    ":4000",
		Handler: routes.GetRoutes(&storeUsecase),
	}
	log.Printf("Runnig server on %s port\n", srv.Addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
