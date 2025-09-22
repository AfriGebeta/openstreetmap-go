package repository

import (
	"gorm.io/gorm"
	"openstreetmap-go/src/db/client"
)

func GetApiSizeLimitByUserId(userId int64) (int64, error) {
	return getApiSizeLimitByUserId(client.DatabaseClients.GetMasterConnection(), userId)
}

func getApiSizeLimitByUserId(client *gorm.DB, userId int64) (int64, error) {
	type Result struct {
		OK int64 `gorm:"column:ok"`
	}

	var result Result
	err := client.Raw("SELECT api_size_limit(?) AS ok", userId).Scan(&result).Error
	if err != nil {
		return 0, err
	}

	return result.OK, nil
}
