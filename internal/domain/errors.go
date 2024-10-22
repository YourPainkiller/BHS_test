package domain

import "errors"

var (
	ErrInvalidUserId       = errors.New("invalid userId")
	ErrInvalidUsername     = errors.New("invalid userName")
	ErrInvalidAssetId      = errors.New("invalid assetId")
	ErrInvalidAssetPrice   = errors.New("invalid assetPrice")
	ErrInvalidAssetName    = errors.New("invalid assetName")
	ErrInvalidAssetDescr   = errors.New("invalid assetDescr")
	ErrInvalidExpiresIn    = errors.New("invalid expiresIn")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUnkown              = errors.New("unkown error")
	ErrUnkownSigningMethod = errors.New("unkown signing method")
	ErrNoSuchUser          = errors.New("no such user")
	ErrNoSuchSession       = errors.New("no such session")
	ErrNoSuchAsset         = errors.New("no such asset")
	ErrAlreadyExists       = errors.New("already exists")
	UniqueErrCode          = "23505"
)

type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
}

type AcceptResponse struct {
	Detail string `json:"detail" example:"success"`
}
