package repository

import (
	"gorm.io/gorm"
	"openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
	"strconv"
	"time"
)

func GetChangeSetById(id int64) (gormModel.Changesets, error) {
	return getChangeSetById(client.DatabaseClients.GetMasterConnection(), id)
}

func CloseChangeSetById(id int64) error {
	return closeChangeSetById(client.DatabaseClients.GetMasterConnection(), id)
}

func InsertIntoChangeset(changeset *gormModel.Changesets, changesetTags []gormModel.ChangesetTags, changesetSubscribers gormModel.ChangesetsSubscribers) (gormModel.Changesets, error) {
	return insertIntoChangeset(client.DatabaseClients.GetMasterConnection(), changeset, changesetTags, changesetSubscribers)
}
func CreateNewNode(currentNode *gormModel.CurrentNodes, nodes *gormModel.Nodes, currentNodeTags []gormModel.CurrentNodeTags, nodeTags []gormModel.NodeTags) error {
	return createNewNode(client.DatabaseClients.GetMasterConnection(), currentNode, nodes, currentNodeTags, nodeTags)
}
func ModifyNode(nodeId string, data map[interface{}]interface{}, nodes *gormModel.Nodes, currentNodeTags []gormModel.CurrentNodeTags, nodeTags []gormModel.NodeTags) error {
	return modifyNode(client.DatabaseClients.GetMasterConnection(), nodeId, data, nodes, currentNodeTags, nodeTags)
}

func modifyNode(client *gorm.DB, nodeId string, data map[interface{}]interface{}, nodes *gormModel.Nodes, currentNodeTags []gormModel.CurrentNodeTags, nodeTags []gormModel.NodeTags) error {
	tx := client.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := client.Model(gormModel.CurrentNodes{}).Where("id = ?", nodeId).Updates(data).Error
	if err != nil {
		return err
	}

	nodeIdInt, err := strconv.ParseInt(nodeId, 10, 64)
	if err != nil {
		return err
	}
	nodes.NodeId = nodeIdInt

	err = client.Model(&gormModel.Nodes{}).Create(&nodes).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, nodeTag := range nodeTags {
		nodeTag.NodeId = nodeIdInt
		err = tx.Model(&gormModel.NodeTags{}).Create(nodeTag).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, currentNodeTag := range currentNodeTags {
		currentNodeTag.NodeId = nodeIdInt
		err = tx.Model(&gormModel.CurrentNodeTags{}).Create(currentNodeTag).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil

}

func UpdateChangeset(id string, data map[interface{}]interface{}) error {
	return updateChangeset(client.DatabaseClients.GetMasterConnection(), id, data)
}

func updateChangeset(client *gorm.DB, id string, data map[interface{}]interface{}) error {
	err := client.Model(gormModel.Changesets{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

func createNewNode(client *gorm.DB, currentNode *gormModel.CurrentNodes, nodes *gormModel.Nodes, currentNodeTags []gormModel.CurrentNodeTags, nodeTags []gormModel.NodeTags) error {
	tx := client.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Model(&gormModel.CurrentNodes{}).Create(&currentNode).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	nodes.NodeId = currentNode.ID
	err = client.Model(&gormModel.Nodes{}).Create(&nodes).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, nodeTag := range nodeTags {
		nodeTag.NodeId = currentNode.ID
		err = tx.Model(&gormModel.NodeTags{}).Create(nodeTag).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, currentNodeTag := range currentNodeTags {
		currentNodeTag.NodeId = currentNode.ID
		err = tx.Model(&gormModel.CurrentNodeTags{}).Create(currentNodeTag).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func insertIntoChangeset(client *gorm.DB, changeset *gormModel.Changesets, changesetTags []gormModel.ChangesetTags, changesetSubscribers gormModel.ChangesetsSubscribers) (gormModel.Changesets, error) {
	err := client.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&gormModel.Changesets{}).Create(changeset).Error
		if err != nil {
			return err
		}

		for _, tag := range changesetTags {
			anTag := tag
			anTag.ChangesetId = changeset.ID
			err := tx.Model(&gormModel.ChangesetTags{}).Create(anTag).Error
			if err != nil {
				return err
			}
		}

		changeSubscribersId := changesetSubscribers
		changeSubscribersId.ChangesetId = changeset.ID
		err = tx.Model(&gormModel.ChangesetsSubscribers{}).Create(changeSubscribersId).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return gormModel.Changesets{}, err
	}
	return *changeset, nil
}

func closeChangeSetById(client *gorm.DB, id int64) error {
	err := client.Model(&gormModel.Changesets{}).
		Where("id = ?", id).
		Update("closed_at", time.Now()).
		Error

	return err
}

func getChangeSetById(client *gorm.DB, id int64) (gormModel.Changesets, error) {
	var changeset gormModel.Changesets

	err := client.Model(&changeset).Where("id = ?", id).First(&changeset).Error
	if err != nil {
		return gormModel.Changesets{}, err
	}
	return changeset, nil
}
