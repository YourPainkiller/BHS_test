package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/YourPainkiller/BHS_test/internal/domain"
	"github.com/YourPainkiller/BHS_test/internal/dto"
	"github.com/YourPainkiller/BHS_test/internal/usecase"
)

type Handlers struct {
	StoreUseCase usecase.StoreUseCase
}

func NewHandlers(storeUseCase usecase.StoreUseCase) *Handlers {
	return &Handlers{
		StoreUseCase: storeUseCase,
	}
}

func (h *Handlers) MainPageForm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("hello"))
}

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req := dto.RegisterUserDto{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error in Register user decoding json: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	ctx := context.Background()
	err := h.StoreUseCase.RegisterUser(ctx, req)

	switch {
	case err != nil:
		switch {
		case errors.Is(err, domain.ErrInvalidUsername):
			http.Error(w, "wrong username", http.StatusUnauthorized)
		default:
			http.Error(w, "unkown error", http.StatusInternalServerError)
		}
	default:
		SendJson(w, map[string]string{"detail": "successed"}, http.StatusAccepted)
	}
}

func SendJson(w http.ResponseWriter, data map[string]string, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Write(jsonResponse)
	return nil
}
