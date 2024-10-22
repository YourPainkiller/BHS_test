package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/YourPainkiller/BHS_test/internal/domain"
	"github.com/YourPainkiller/BHS_test/internal/dto"
	"github.com/YourPainkiller/BHS_test/internal/usecase"
)

const ACCESSSECRET = "access"
const REFRESHSECRET = "refresh"

type Handlers struct {
	StoreUseCase usecase.StoreUseCase
}

func NewHandlers(storeUseCase usecase.StoreUseCase) *Handlers {
	return &Handlers{
		StoreUseCase: storeUseCase,
	}
}

func (h *Handlers) MainPageForm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello"))
}

// Register user godoc
// @Summary		Register user
// @Description	With this command you can register user
// @Accept			json
// @Produce		json
// @Param UserData body dto.UserCredentials true "your credentils"
// @Tags auth
// @Success		200	{object}	domain.AcceptResponse
// @Failure		400	{object}	domain.ErrorResponse
// @Router			/api/auth/register [post]
func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req := dto.RegisterUserDto{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error in Register user decoding json: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	req.UserPassword = genHash(req.UserPassword)
	user, err := domain.NewUser(req.UserName, req.UserPassword)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidUsername):
			http.Error(w, "invalid username", http.StatusBadRequest)
			return
		default:
			log.Printf("Register user new user: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
	}

	ctx := context.Background()
	err = h.StoreUseCase.RegisterUser(ctx, user.ToDTO())

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAlreadyExists):
			http.Error(w, "username already taken", http.StatusBadRequest)
			return
		default:
			log.Println(err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
	}
	SendJson(w, map[string]string{"detail": "successed"}, http.StatusOK)
}

// Login user godoc
// @Summary		Login user
// @Description	With this command you can login user
// @Accept			json
// @Produce		json
// @Param UserData body dto.UserCredentials true "your credentils"
// @Tags auth
// @Success		200	{object}	domain.AcceptResponse
// @Failure		400	{object}	domain.ErrorResponse
// @Router			/api/auth/login [post]
func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	//TODO принимать с фронта fingerpring браузера
	req := dto.RegisterUserDto{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error in Register user decoding json: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	req.UserPassword = genHash(req.UserPassword)

	user, err := domain.NewUser(req.UserName, req.UserPassword)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidUsername):
			http.Error(w, "invalid username", http.StatusBadRequest)
			return
		default:
			log.Printf("Login user new user: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
	}

	ctx := context.Background()

	userId, err := h.StoreUseCase.LoginUser(ctx, user.ToDTO())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoSuchUser):
			http.Error(w, "no such user or wrong password", http.StatusBadRequest)
			return
		case errors.Is(err, domain.ErrInvalidCredentials):
			http.Error(w, "no such user or wrong password", http.StatusBadRequest)
			return
		default:
			log.Printf("error in login: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
	}

	session := dto.RefreshSessionDto{UserId: userId}
	session.RefreshToken, session.Expires, err = genToken(userId, req.UserName, 7200, REFRESHSECRET)
	if err != nil {
		log.Printf("error in creating refresh token: %v\n", err)
		http.Error(w, "unkown error", http.StatusInternalServerError)
		return
	}
	session.Ip = readUserIP(r)
	session.Fingerprint = "somedata"
	session.CreatedAt = time.Now()

	//TODO добавить проверку на количество сессий
	err = h.StoreUseCase.SetSession(ctx, session)
	if err != nil {
		log.Printf("error in setting session: %v\n", err)
		http.Error(w, "unkown error", http.StatusInternalServerError)
		return
	}

	accessToken, exp, err := genToken(userId, req.UserName, 60, ACCESSSECRET)
	if err != nil {
		log.Printf("error in creating access token: %v\n", err)
		http.Error(w, "unkown error", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:    "accessToken",
		Value:   accessToken,
		Path:    "/api/auth",
		Expires: exp,
	}
	http.SetCookie(w, cookie)

	cookie = &http.Cookie{
		Name:    "refreshToken",
		Value:   session.RefreshToken,
		Path:    "/api/auth",
		Expires: session.Expires,
	}
	http.SetCookie(w, cookie)
	SendJson(w, map[string]string{"detail": "successed"}, http.StatusOK)
}

// Adding asset godoc
// @Summary		Add asset
// @Description	Adding asset to store
// @Accept			json
// @Produce		json
// @Param AssetData body dto.Add true "asset info"
// @Tags asset
// @Success		200	{object}	domain.AcceptResponse
// @Failure		400	{object}	domain.ErrorResponse
// @Router			/api/auth/add [post]
func (h *Handlers) AddAsset(w http.ResponseWriter, r *http.Request) {
	userId, _, _, err := parseCookie(r, "accessToken")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "missing cookie", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrInvalidCredentials):
			http.Error(w, "wrong cookie", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrInvalidExpiresIn):
			http.Error(w, "cookie expired", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrUnkownSigningMethod):
			http.Error(w, "unknow token", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrUnkown):
			log.Printf("error in adding asset: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		default:
			log.Printf("error in adding asset: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
	}

	req := dto.AddAssetDto{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error in add asset decoding incoming json: %v\n", err)
		http.Error(w, "unkown error", http.StatusInternalServerError)
		return
	}

	asset, err := domain.NewAsset(userId, req.AssetPrice, req.AssetName, req.AssetDescr)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidAssetName):
			http.Error(w, "wrong asset name", http.StatusBadRequest)
			return
		case errors.Is(err, domain.ErrInvalidAssetPrice):
			http.Error(w, "wrong asset price", http.StatusBadRequest)
			return
		case errors.Is(err, domain.ErrInvalidAssetDescr):
			http.Error(w, "wrong asset description", http.StatusBadRequest)
			return
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	err = h.StoreUseCase.AddAsset(context.Background(), asset.ToDTO())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAlreadyExists):
			http.Error(w, "asset already exists", http.StatusBadRequest)
			return
		default:
			log.Printf("error in adding asset: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
	}
	SendJson(w, map[string]string{"detail": "successed"}, http.StatusOK)
}

// Deleting asset godoc
// @Summary		Delete asset
// @Description	Deleting asset from store
// @Accept			json
// @Produce		json
// @Param AssetData body dto.Delete true "asset name"
// @Tags asset
// @Success		200	{object}	domain.AcceptResponse
// @Failure		400	{object}	domain.ErrorResponse
// @Router			/api/auth/delete [post]
func (h *Handlers) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	userId, _, _, err := parseCookie(r, "accessToken")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "missing cookie", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrInvalidCredentials):
			http.Error(w, "wrong cookie", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrInvalidExpiresIn):
			http.Error(w, "cookie expired", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrUnkownSigningMethod):
			http.Error(w, "unknow token", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrUnkown):
			log.Printf("error in deleting asset: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		default:
			log.Printf("error in deleting asset: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
	}

	req := dto.DeleteAssetDto{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error in add asset decoding incoming json: %v\n", err)
		http.Error(w, "unkown error", http.StatusInternalServerError)
		return
	}
	req.UserId = userId

	err = h.StoreUseCase.DeleteAsset(context.Background(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoSuchAsset):
			http.Error(w, "no such asset", http.StatusBadRequest)
			return
		default:
			log.Printf("error in deleting asset: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
	}
	SendJson(w, map[string]string{"detail": "successed"}, http.StatusOK)
}

// Buying asset godoc
// @Summary		Buy asset
// @Description	Buying asset from store
// @Accept			json
// @Produce		json
// @Param AssetData body dto.Buy true "asset name and count"
// @Tags asset
// @Success		200	{object}	domain.AcceptResponse
// @Failure		400	{object}	domain.ErrorResponse
// @Router			/api/auth/buy [post]
func (h *Handlers) BuyAsset(w http.ResponseWriter, r *http.Request) {
	userId, _, _, err := parseCookie(r, "accessToken")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "missing cookie", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrInvalidCredentials):
			http.Error(w, "wrong cookie", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrInvalidExpiresIn):
			http.Error(w, "cookie expired", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrUnkownSigningMethod):
			http.Error(w, "unknow token", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrUnkown):
			log.Printf("error in deleting asset: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		default:
			log.Printf("error in deleting asset: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
	}

	req := dto.BuyAssetDto{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error in buy asset decoding incoming json: %v\n", err)
		http.Error(w, "unkown error", http.StatusInternalServerError)
		return
	}
	req.UserId = userId
	asset, err := h.StoreUseCase.BuyAsset(context.Background(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoSuchAsset):
			http.Error(w, "no such asset", http.StatusBadRequest)
			return
		default:
			log.Printf("error in buying asset: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
	}

	log.Printf("user=%d bought item=%d. Total price = %d\n", req.UserId, asset.AssetId, asset.AssetPrice*req.Count)
	SendJson(w, map[string]string{"detail": "successed"}, http.StatusOK)
}

// Refresh cookie godoc
// @Summary		Refresh cookie
// @Description	Refresh your auth
// @Produce		json
// @Tags auth
// @Success		200	{object}	domain.AcceptResponse
// @Failure		400	{object}	domain.ErrorResponse
// @Router			/api/auth/refresh [get]
func (h *Handlers) RefreshSession(w http.ResponseWriter, r *http.Request) {
	userId, username, refreshToken, err := parseCookie(r, "refreshToken")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "missing cookie", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrInvalidCredentials):
			http.Error(w, "wrong cookie", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrInvalidExpiresIn):
			http.Error(w, "cookie expired", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrUnkownSigningMethod):
			http.Error(w, "unknow token", http.StatusUnauthorized)
			return
		case errors.Is(err, domain.ErrUnkown):
			log.Printf("error in deleting asset: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		default:
			log.Printf("error in deleting asset: %v\n", err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
	}

	req := dto.UpdateRefreshDto{UserId: userId, PriviousRefresh: refreshToken}
	req.RefreshToken, req.Expires, err = genToken(userId, username, 7200, REFRESHSECRET)
	if err != nil {
		log.Printf("error in creating refresh token: %v\n", err)
		http.Error(w, "unkown error", http.StatusInternalServerError)
		return
	}
	accessToken, exp, err := genToken(userId, username, 60, ACCESSSECRET)
	if err != nil {
		log.Printf("error in creating access token: %v\n", err)
		http.Error(w, "unkown error", http.StatusInternalServerError)
		return
	}

	req.Ip = readUserIP(r)
	req.Fingerprint = "somedata"
	req.CreatedAt = time.Now()

	err = h.StoreUseCase.Refresh(context.Background(), req)
	if err != nil {
		log.Printf("error in refreshing %v", err)
		http.Error(w, "unkown error", http.StatusInternalServerError)
	}

	cookie := &http.Cookie{
		Name:    "accessToken",
		Value:   accessToken,
		Path:    "/api/auth",
		Expires: exp,
	}
	http.SetCookie(w, cookie)

	cookie = &http.Cookie{
		Name:    "refreshToken",
		Value:   req.RefreshToken,
		Path:    "/api/auth",
		Expires: req.Expires,
	}
	http.SetCookie(w, cookie)

	SendJson(w, map[string]string{"detail": "successed"}, http.StatusOK)
}
