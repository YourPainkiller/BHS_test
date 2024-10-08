package dto

type RefreshSessionDto struct {
	UserId       int    `json:"userId" db:"user_id"`
	RefreshToken string `json:"refreshToken" db:"refresh_token"`
	Fingerprint  string `json:"fingerprint" db:"fingerprint"`
	Ip           string `json:"ip" db:"ip"`
	ExpiresIn    int    `json:"expiresIn" db:"expires_in"`
	CreatedAt    string `json:"createdAt" db:"created_at"`
}
