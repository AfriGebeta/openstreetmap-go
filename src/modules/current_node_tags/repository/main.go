package repository

import (
	"gorm.io/gorm"
	dbClient "openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
)

func DeleteCurrentNodeTags(id string) error {
	return deleteCurrentNodeTags(dbClient.DatabaseClients.GetMasterConnection(), id)
}

func deleteCurrentNodeTags(client *gorm.DB, id string) error {
	if err := client.Where("node_id = ?", id).Delete(&gormModel.NodeTags{}).Error; err != nil {
		return err
	}
	return nil
}
