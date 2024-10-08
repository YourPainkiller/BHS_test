package domain

import "errors"

var (
	ErrInvalidUserId     = errors.New("invalid userId")
	ErrInvalidUsername   = errors.New("invalid userName")
	ErrInvalidAssetId    = errors.New("invalid assetId")
	ErrInvalidAssetPrice = errors.New("invalid assetPrice")
	ErrInvalidAssetName  = errors.New("invalid assetName")
	ErrInvalidAssetDescr = errors.New("invalid assetDescr")
	ErrInvalidExpiresIn  = errors.New("invalid expiresIn")
)
