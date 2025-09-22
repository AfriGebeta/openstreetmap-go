package respository

import (
	"gorm.io/gorm"
	dbClient "openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
	db "openstreetmap-go/src/db/generated/prisma"
)

func InsertUserRole(userrole db.InnerUserRoles) error {
	return insertUserRole(dbClient.DatabaseClients.GetMasterConnection(), userrole)
}
func GetUserRoleByUserId(userId int64) ([]gormModel.UserRoles, error) {
	return getUserRoleByUserId(dbClient.DatabaseClients.GetMasterConnection(), userId)
}

func getUserRoleByUserId(client *gorm.DB, userId int64) ([]gormModel.UserRoles, error) {
	var roles []gormModel.UserRoles
	err := client.Model(&roles).Where("user_id = ?", userId).Scan(&roles).Error

	if err != nil {
		return nil, err
	}
	return roles, nil
}

func insertUserRole(client *gorm.DB, userrole db.InnerUserRoles) error {
	err := client.Create(&userrole).Error
	if err != nil {
		return err
	}
	return nil
}
