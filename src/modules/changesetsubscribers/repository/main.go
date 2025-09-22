package repository

import (
	"gorm.io/gorm"
	"openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
)

func InsertChangesetSubscriber(changesetSubscriber *gormModel.ChangesetsSubscribers) error {
	return insertChangesetSubscriber(client.DatabaseClients.GetMasterConnection(), changesetSubscriber)
}

func insertChangesetSubscriber(client *gorm.DB, changesetSubscriber *gormModel.ChangesetsSubscribers) error {
	err := client.Model(&gormModel.ChangesetsSubscribers{}).Create(changesetSubscriber).Error
	if err != nil {
		return err
	}
	return nil
}
