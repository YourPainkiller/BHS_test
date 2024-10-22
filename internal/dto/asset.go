package dto

type AssetDto struct {
	AssetId    int    `json:"assetId" db:"id"`
	UserId     int    `json:"userId" db:"user_id"`
	AssetName  string `json:"assetName" db:"name"`
	AssetDescr string `json:"assetDescr" db:"descr"`
	AssetPrice int    `json:"assetPrice" db:"price"`
}

type AddAssetDto struct {
	AssetName  string `json:"assetName" db:"name"`
	AssetDescr string `json:"assetDescr" db:"descr"`
	AssetPrice int    `json:"assetPrice" db:"price"`
}

type DeleteAssetDto struct {
	UserId    int    `json:"userId" db:"user_id"`
	AssetName string `json:"assetName" db:"name"`
}

type BuyAssetDto struct {
	UserId    int    `json:"userId" db:"user_id"`
	AssetName string `json:"assetName" db:"name"`
	Count     int    `json:"count"`
}

// models for swagger
type Add struct {
	AssetName  string `json:"assetName" example:"tree"`
	AssetDescr string `json:"assetDescr" example:"lorem ipsum"`
	AssetPrice int    `json:"assetPrice" example:"10"`
}

type Delete struct {
	AssetName string `json:"assetName" example:"tree"`
}

type Buy struct {
	AssetName string `json:"assetName" example:"tree"`
	Count     int    `json:"count" example:"3"`
}
