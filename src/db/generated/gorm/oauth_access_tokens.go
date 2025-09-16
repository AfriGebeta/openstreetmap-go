package models

import "time"

type OauthAccessTokens struct {
	ID                   int64             `gorm:"primaryKey;column:id" json:"id"`
	ResourceOwnerId      *int64            `gorm:"column:resource_owner_id" json:"resource_owner_id"`
	ApplicationId        int64             `gorm:"column:application_id" json:"application_id"`
	ExpiresIn            *int              `gorm:"column:expires_in" json:"expires_in"`
	RevokedAt            *time.Time        `gorm:"column:revoked_at" json:"revoked_at"`
	CreatedAt            time.Time         `gorm:"column:created_at" json:"created_at"`
	Scopes               *string           `gorm:"column:scopes" json:"scopes"`
	PreviousRefreshToken string            `gorm:"column:previous_refresh_token" json:"previous_refresh_token"`
	Token                string            `gorm:"column:token" json:"token"`
	RefreshToken         *string           `gorm:"column:refresh_token" json:"refresh_token"`
	Application          OauthApplications `gorm:"foreignKey:ApplicationId;references:ID"`
	ResourceOwner        Users             `gorm:"foreignKey:ResourceOwnerId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *OauthAccessTokens) TableName() string {
	return "oauth_access_tokens"
}
