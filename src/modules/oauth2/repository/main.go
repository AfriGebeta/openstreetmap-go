package repository

import (
	"gorm.io/gorm"
	"openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
)

func GetOAuthAccessToken(token string) (gormModel.OauthAccessTokens, error) {
	return getOauthAccessToken(client.DatabaseClients.GetMasterConnection(), token)
}

func getOauthAccessToken(client *gorm.DB, token string) (gormModel.OauthAccessTokens, error) {
	var OauthAccessTokens gormModel.OauthAccessTokens
	err := client.Model(&OauthAccessTokens).Where("token = ?", token).First(&OauthAccessTokens).Error
	if err != nil {
		return gormModel.OauthAccessTokens{}, err
	}
	return OauthAccessTokens, nil
}
