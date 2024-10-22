package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/YourPainkiller/BHS_test/internal/domain"
	"github.com/golang-jwt/jwt/v4"
)

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

func genToken(userId int, username string, timedur int, secret string) (string, time.Time, error) {
	hmacSampleSecret := []byte(secret)
	exp := time.Now().Add(time.Minute * time.Duration(timedur))
	u := strconv.Itoa(userId)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expired":  exp.Format(time.DateTime),
		"userId":   u,
		"username": username,
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", time.Now(), err
	}
	return tokenString, exp, nil
}

func parseCookie(r *http.Request, key string) (int, string, string, error) {
	cookie, err := r.Cookie(key)
	if err != nil {
		return 0, "", "", err
	}
	if cookie.Value == "" {
		return 0, "", "", http.ErrNoCookie
	}

	hmacSampleSecret := []byte(ACCESSSECRET)
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", domain.ErrUnkownSigningMethod, token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil && errors.Is(err, domain.ErrUnkownSigningMethod) {
		return 0, "", "", err
	}
	if !token.Valid {
		return 0, "", "", domain.ErrInvalidCredentials
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		username := claims["username"].(string)
		e, _ := claims["expired"].(string)
		exp, _ := time.Parse(time.DateTime, e)
		u, _ := claims["userId"].(string)
		userId, _ := strconv.Atoi(u)

		switch {
		case time.Now().After(exp):
			return 0, "", "", domain.ErrInvalidExpiresIn
		default:
			return userId, username, cookie.Value, nil
		}
	}
	return 0, "", "", domain.ErrUnkown
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
