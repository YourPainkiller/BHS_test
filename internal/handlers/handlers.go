package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/YourPainkiller/BHS_test/internal/domain"
	"github.com/YourPainkiller/BHS_test/internal/dto"
	"github.com/YourPainkiller/BHS_test/internal/repository/postgres"
	"github.com/YourPainkiller/BHS_test/internal/usecase"
	"github.com/golang-jwt/jwt/v4"
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
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("hello"))
}

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req := dto.UserDto{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error in Register user decoding json: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	req.UserPassword = genHash(req.UserPassword)

	ctx := context.Background()
	err := h.StoreUseCase.RegisterUser(ctx, req)

	switch {
	case err != nil:
		switch {
		case errors.Is(err, domain.ErrInvalidUsername):
			http.Error(w, "wrong username", http.StatusUnauthorized)
		case postgres.UnwrapPgCode(err) == domain.ErrSameUniqeCode:
			http.Error(w, "username already taken", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
		}
	default:
		SendJson(w, map[string]string{"detail": "successed"}, http.StatusAccepted)
	}
}

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	//TODO принимать с фронта fingerpring браузера
	req := dto.UserDto{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error in Register user decoding json: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	req.UserPassword = genHash(req.UserPassword)

	ctx := context.Background()
	userId, err := h.StoreUseCase.LoginUser(ctx, req)

	switch {
	case err != nil:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		session := dto.RefreshSessionDto{UserId: userId}
		session.RefreshToken, err = genToken(req.UserName, REFRESHSECRET)
		if err != nil {
			log.Println(err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}
		session.Ip = readUserIP(r)
		session.Fingerprint = "somedata"
		session.ExpiresIn = 86400
		session.CreatedAt = time.Now().Format(time.DateTime)

		//TODO добавить проверку на количество сессий
		err = h.StoreUseCase.SetSession(ctx, session)
		if err != nil {
			log.Println(err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}

		accessToken, err := genToken(req.UserName, ACCESSSECRET)
		if err != nil {
			log.Println(err)
			http.Error(w, "unkown error", http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{
			Name:   "accessToken",
			Value:  accessToken,
			Path:   "/api/auth",
			MaxAge: 3600,
		}
		http.SetCookie(w, cookie)
		cookie = &http.Cookie{
			Name:   "refreshToken",
			Value:  session.RefreshToken,
			Path:   "/api/auth",
			MaxAge: session.ExpiresIn,
		}
		http.SetCookie(w, cookie)
		SendJson(w, map[string]string{"detail": "successed"}, http.StatusAccepted)
	}
}

func (h *Handlers) AddAsset(w http.ResponseWriter, r *http.Request) {
	username, err := parseCookie(r)
	switch {
	case errors.Is(err, http.ErrNoCookie):
		http.Error(w, "missing cookie", http.StatusForbidden)
		return
	case errors.Is(err, domain.ErrInvalidCredentials):
		http.Error(w, "wrong cookie", http.StatusForbidden)
		return
	default:
		log.Println(err)
		http.Error(w, "unkown error", http.StatusInternalServerError)
		return
	}
	fmt.Println(username)

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

func genHash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	//TODO: Передавать секрет вместо nil
	fh := h.Sum(nil)
	temp := hex.EncodeToString(fh)
	return string(temp)
}

func genToken(username string, secret string) (string, error) {
	hmacSampleSecret := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func parseCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("accessToken")
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", http.ErrNoCookie
	}

	hmacSampleSecret := []byte(ACCESSSECRET)
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", domain.ErrInvalidCredentials
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		username, _ := claims["username"].(string)
		return username, nil
	}
	return "", nil
}
