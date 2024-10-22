package dto

import "time"

type RefreshSessionDto struct {
	UserId       int       `json:"userId" db:"user_id"`
	RefreshToken string    `json:"refreshToken" db:"refresh_token"`
	Fingerprint  string    `json:"fingerprint" db:"fingerprint"`
	Ip           string    `json:"-" db:"ip"`
	Expires      time.Time `json:"expires" db:"expires_in"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type UpdateRefreshDto struct {
	UserId          int    `json:"userId" db:"user_id"`
	RefreshToken    string `json:"refreshToken" db:"refresh_token"`
	PriviousRefresh string
	Fingerprint     string    `json:"fingerprint" db:"fingerprint"`
	Ip              string    `json:"-" db:"ip"`
	Expires         time.Time `json:"expires" db:"expires_in"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}
