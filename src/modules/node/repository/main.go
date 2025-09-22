package repository

import (
	"gorm.io/gorm"
	dbClient "openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
)

func CreateNode(node *gormModel.Nodes) error {
	return createNode(dbClient.DatabaseClients.GetMasterConnection(), node)
}

func createNode(client *gorm.DB, currentNode *gormModel.Nodes) error {
	err := client.Model(&gormModel.Nodes{}).Create(&currentNode).Error
	if err != nil {
		return err
	}
	return nil
}
