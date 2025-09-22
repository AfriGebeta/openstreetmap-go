package repository

import (
	"gorm.io/gorm"
	dbClient "openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
	"time"
)

func GetBlockedUser(userId int64, time time.Time) (gormModel.UserBlocks, error) {
	return getBlockedUser(dbClient.DatabaseClients.GetMasterConnection(), userId, time)
}

func getBlockedUser(client *gorm.DB, userId int64, time time.Time) (gormModel.UserBlocks, error) {
	var users gormModel.UserBlocks

	err := client.Model(&users).
		Where("user_id = ? AND (needs_view = TRUE OR ends_at > ?)", userId, time).
		First(&users).Error

	if err != nil {
		return gormModel.UserBlocks{}, err
	}

	return users, nil
}
