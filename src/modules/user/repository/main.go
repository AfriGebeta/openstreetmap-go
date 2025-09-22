package repository

import (
	"gorm.io/gorm"
	dbClient "openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
	"time"
)

func InsertUser(users *gormModel.Users) error {
	return insertUser(dbClient.DatabaseClients.GetMasterConnection(), users)
}

func GetUser(condition map[string]interface{}) (gormModel.Users, error) {
	return getUser(dbClient.DatabaseClients.GetRandomConnection(), condition)
}

func UpdateUser(data map[string]interface{}, condition map[string]interface{}) (gormModel.Users, error) {
	return updateUser(dbClient.DatabaseClients.GetMasterConnection(), data, condition)
}

func GetUserByEmailorDisplayName(email, displayName string) (gormModel.Users, error) {
	return getUserByEmailOrDisplayName(dbClient.DatabaseClients.GetMasterConnection(), email, displayName)
}

func GetBlockedUser(userId string, time time.Time) (gormModel.Users, error) {
	return getBlockedUser(dbClient.DatabaseClients.GetMasterConnection(), userId, time)
}

func UpdateChangesetCount(userId string) error {
	return updateChangesetCount(dbClient.DatabaseClients.GetMasterConnection(), userId)
}

func updateChangesetCount(client *gorm.DB, userId string) error {
	err := client.Model(&gormModel.Users{}).
		Where("id = ?", userId).
		Update("changesets_count", gorm.Expr("COALESCE(changesets_count, 0) + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func getBlockedUser(client *gorm.DB, userId string, time time.Time) (gormModel.Users, error) {
	var users gormModel.Users

	err := client.Model(&users).
		Where("user_id = ? AND (needs_view = TRUE OR ends_at > ?)", userId, time).
		First(&users).Error

	if err != nil {
		return gormModel.Users{}, err
	}

	return users, nil
}

func getUserByEmailOrDisplayName(client *gorm.DB, email string, displayName string) (gormModel.Users, error) {
	var user gormModel.Users
	err := client.Where("LOWER(email) = ?", email).
		Or("display_name ~* ?", displayName).
		First(&user).Error
	if err != nil {
		return gormModel.Users{}, err
	}
	return user, nil
}

func getUser(client *gorm.DB, condition map[string]interface{}) (gormModel.Users, error) {
	var user gormModel.Users
	err := client.Where(condition).First(&user).Error
	if err != nil {
		return gormModel.Users{}, err
	}
	return user, nil
}

func insertUser(client *gorm.DB, users *gormModel.Users) error {

	err := client.Create(&users).Error
	if err != nil {
		return err
	}

	return nil
}

func updateUser(client *gorm.DB, data map[string]interface{}, condition map[string]interface{}) (gormModel.Users, error) {
	var user gormModel.Users
	err := client.Model(&user).Where(condition).Updates(data).First(&user).Error
	if err != nil {
		return gormModel.Users{}, err
	}
	return user, nil
}
