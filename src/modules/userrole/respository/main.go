package respository

import (
	"gorm.io/gorm"
	dbClient "openstreetmap-go/src/db/client"
	db "openstreetmap-go/src/db/generated/prisma"
)

func InsertUserRole(userrole db.InnerUserRoles) error {
	return insertUserRole(dbClient.DatabaseClients.GetMasterConnection(), userrole)
}

func insertUserRole(client *gorm.DB, userrole db.InnerUserRoles) error {
	err := client.Create(&userrole).Error
	if err != nil {
		return err
	}
	return nil
}
