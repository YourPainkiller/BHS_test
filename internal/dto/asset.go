package dto

type AssetDto struct {
	AssetId    int    `json:"assetId" db:"id"`
	UserId     int    `json:"userId" db:"user_id"`
	AssetName  string `json:"assetName" db:"name"`
	AssetDescr string `json:"assetDescr" db:"descr"`
	AssetPrice int    `json:"assetPrice" db:"price"`
}
