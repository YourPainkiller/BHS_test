package domain

import (
	"time"

	"github.com/YourPainkiller/BHS_test/internal/dto"
)

type RefreshSession struct {
	userId       int
	refreshToken string
	fingerprint  string
	ip           string
	expiresIn    time.Time
	createdAt    time.Time
}

func NewRefreshSession(userId int, expiresIn, createdAt time.Time, refreshToken, fingerprint, ip string) (*RefreshSession, error) {
	rs := RefreshSession{}
	err := rs.SetUserId(userId)
	if err != nil {
		return nil, err
	}

	err = rs.SetExpiresIn(expiresIn)
	if err != nil {
		return nil, err
	}

	err = rs.SetRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	err = rs.Setfingerprint(fingerprint)
	if err != nil {
		return nil, err
	}

	err = rs.SetIp(ip)
	if err != nil {
		return nil, err
	}

	err = rs.SetCreatedAt(createdAt)
	if err != nil {
		return nil, err
	}

	return &rs, nil

}

func (rs *RefreshSession) SetUserId(userId int) error {
	if userId < 1 {
		return ErrInvalidUserId
	}
	rs.userId = userId
	return nil
}

func (rs *RefreshSession) SetExpiresIn(expiresIn time.Time) error {
	rs.expiresIn = expiresIn
	return nil
}

func (rs *RefreshSession) SetRefreshToken(refreshToken string) error {
	rs.refreshToken = refreshToken
	return nil
}

func (rs *RefreshSession) Setfingerprint(fingerprint string) error {
	rs.fingerprint = fingerprint
	return nil
}

func (rs *RefreshSession) SetIp(ip string) error {
	rs.ip = ip
	return nil
}

func (rs *RefreshSession) SetCreatedAt(createdAt time.Time) error {
	rs.createdAt = createdAt
	return nil
}

func (rs *RefreshSession) ToDTO() dto.RefreshSessionDto {
	return dto.RefreshSessionDto{
		UserId:       rs.userId,
		RefreshToken: rs.refreshToken,
		Fingerprint:  rs.fingerprint,
		Ip:           rs.ip,
		Expires:      rs.expiresIn,
		CreatedAt:    rs.createdAt,
	}
}
