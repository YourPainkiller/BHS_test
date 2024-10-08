package domain

import "time"

type RefreshSession struct {
	userId       int
	refreshToken string
	fingerprint  string
	ip           string
	expiresIn    int
	createdAt    time.Time
}

func NewRefreshSession(userId, expiresIn int, refreshToken, fingerprint, ip string, createdAt time.Time) (*RefreshSession, error) {
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

func (rs *RefreshSession) SetExpiresIn(expiresIn int) error {
	if expiresIn < 1 {
		return ErrInvalidExpiresIn
	}
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
