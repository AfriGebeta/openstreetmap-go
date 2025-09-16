package models

import "time"

type OauthApplications struct {
	ID                int64               `gorm:"primaryKey;column:id" json:"id"`
	OwnerId           int64               `gorm:"column:owner_id" json:"owner_id"`
	Confidential      bool                `gorm:"column:confidential" json:"confidential"`
	CreatedAt         time.Time           `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time           `gorm:"column:updated_at" json:"updated_at"`
	OwnerType         string              `gorm:"column:owner_type" json:"owner_type"`
	Name              string              `gorm:"column:name" json:"name"`
	Uid               string              `gorm:"column:uid" json:"uid"`
	Secret            string              `gorm:"column:secret" json:"secret"`
	RedirectUri       string              `gorm:"column:redirect_uri" json:"redirect_uri"`
	Scopes            string              `gorm:"column:scopes" json:"scopes"`
	Owner             Users               `gorm:"foreignKey:OwnerId;references:ID"`
	OauthAccessGrants []OauthAccessGrants `gorm:"foreignKey:ApplicationId"`
	OauthAccessTokens []OauthAccessTokens `gorm:"foreignKey:ApplicationId"`
}

// TableName sets the insert table name for this struct type
func (m *OauthApplications) TableName() string {
	return "oauth_applications"
}
