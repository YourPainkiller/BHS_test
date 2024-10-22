package domain

import "github.com/YourPainkiller/BHS_test/internal/dto"

type Asset struct {
	assetId    int
	userId     int
	assetName  string
	assetDescr string
	assetPrice int
}

const MAXASSETNAMELEN = 12

func NewAsset(userId, assetPrice int, assetName, assetDescr string) (*Asset, error) {
	asset := Asset{}
	err := asset.SetAssetId()
	if err != nil {
		return nil, err
	}

	err = asset.SetUserId(userId)
	if err != nil {
		return nil, err
	}

	err = asset.SetAssetPrice(assetPrice)
	if err != nil {
		return nil, err
	}

	err = asset.SetAssetName(assetName)
	if err != nil {
		return nil, err
	}

	err = asset.SetAssetDescr(assetDescr)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (a *Asset) SetAssetId() error {
	a.assetId = 0
	return nil
}

func (a *Asset) SetUserId(userId int) error {
	if userId < 1 {
		return ErrInvalidUserId
	}
	a.userId = userId
	return nil
}

func (a *Asset) SetAssetPrice(assetPrice int) error {
	if assetPrice < -1 {
		return ErrInvalidAssetPrice
	}
	a.assetPrice = assetPrice
	return nil
}

func (a *Asset) SetAssetName(assetName string) error {
	if len(assetName) < 1 || len(assetName) > MAXASSETNAMELEN {
		return ErrInvalidAssetName
	}
	a.assetName = assetName
	return nil
}

func (a *Asset) SetAssetDescr(assetDescr string) error {
	if len(assetDescr) > 10000 {
		return ErrInvalidAssetDescr
	}
	a.assetDescr = assetDescr
	return nil
}

func (a *Asset) ToDTO() dto.AssetDto {
	return dto.AssetDto{
		AssetId:    a.assetId,
		UserId:     a.userId,
		AssetName:  a.assetName,
		AssetDescr: a.assetDescr,
		AssetPrice: a.assetPrice,
	}
}
