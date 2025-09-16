package aclRespository

import (
	"gorm.io/gorm"
	dbClient "openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
)

func AclCheck(domains, k []string) ([]gormModel.Acls, error) {
	conn := dbClient.DatabaseClients.GetRandomConnection()
	return aclCheck(conn, domains, k)
}

func aclCheck(client *gorm.DB, domains, k []string) ([]gormModel.Acls, error) {
	var result []gormModel.Acls
	// check client
	err := client.
		Where("domain IN (?)", domains).
		Where("k IN (?)", k).
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
