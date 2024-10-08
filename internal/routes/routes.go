package routes

import (
	"net/http"

	"github.com/YourPainkiller/BHS_test/internal/handlers"
	"github.com/YourPainkiller/BHS_test/internal/usecase"
)

func GetRoutes(storeUseCase *usecase.StoreUseCase) *http.ServeMux {
	h := handlers.NewHandlers(*storeUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.MainPageForm)
	mux.HandleFunc("POST /api/auth/register", h.RegisterUser)
	return mux
}
