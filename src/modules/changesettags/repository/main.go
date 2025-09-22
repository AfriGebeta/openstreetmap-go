package repository

import (
	"gorm.io/gorm"
	"openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
)

func CreatechangesetTags(tags []gormModel.ChangesetTags) error {
	return createChangesetTags(client.DatabaseClients.GetMasterConnection(), tags)
}

func createChangesetTags(client *gorm.DB, tags []gormModel.ChangesetTags) error {
	for _, tag := range tags {
		err := client.Model(&gormModel.ChangesetTags{}).Create(tag).Error
		if err != nil {
			return err
		}
	}

	return nil
}
