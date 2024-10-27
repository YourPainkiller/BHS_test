package routes

import (
	"net/http"
	"os"

	"github.com/YourPainkiller/BHS_test/internal/handlers"
	"github.com/YourPainkiller/BHS_test/internal/usecase"
)

func GetRoutes(storeUseCase *usecase.StoreUseCase) *http.ServeMux {
	h := handlers.NewHandlers(*storeUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.MainPageForm)
	mux.HandleFunc("POST /api/auth/register", h.RegisterUser)
	mux.HandleFunc("POST /api/auth/login", h.LoginUser)
	mux.HandleFunc("POST /api/auth/add", h.AddAsset)
	mux.HandleFunc("POST /api/auth/delete", h.DeleteAsset)
	mux.HandleFunc("POST /api/auth/buy", h.BuyAsset)
	mux.HandleFunc("GET /api/auth/refresh", h.RefreshSession)
	mux.HandleFunc("GET /api/auth/logout", h.Logout)
	mux.HandleFunc("GET /swagger-ui", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		b, _ := os.ReadFile("./internal/service/static/index.html")
		w.Write(b)
	})
	mux.HandleFunc("GET /swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, _ := os.ReadFile("./docs/swagger.json")
		w.Write(b)
	})
	return mux
}
