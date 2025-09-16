package repository

import (
	"gorm.io/gorm"
	dbClient "openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
)

func InsertUser(users gormModel.Users) (gormModel.Users, error) {
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

func insertUser(client *gorm.DB, users gormModel.Users) (gormModel.Users, error) {

	err := client.Create(&users).Error
	if err != nil {
		return gormModel.Users{}, err
	}

	return users, nil
}

func updateUser(client *gorm.DB, data map[string]interface{}, condition map[string]interface{}) (gormModel.Users, error) {
	var user gormModel.Users
	err := client.Updates(data).Where(condition).First(&user).Error
	if err != nil {
		return gormModel.Users{}, err
	}
	return user, nil
}
