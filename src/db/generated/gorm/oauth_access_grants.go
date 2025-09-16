package models

import "time"

type OauthAccessGrants struct {
	ID                  int64                 `gorm:"primaryKey;column:id" json:"id"`
	ResourceOwnerId     int64                 `gorm:"column:resource_owner_id" json:"resource_owner_id"`
	ApplicationId       int64                 `gorm:"column:application_id" json:"application_id"`
	ExpiresIn           int                   `gorm:"column:expires_in" json:"expires_in"`
	CreatedAt           time.Time             `gorm:"column:created_at" json:"created_at"`
	RevokedAt           *time.Time            `gorm:"column:revoked_at" json:"revoked_at"`
	Token               string                `gorm:"column:token" json:"token"`
	RedirectUri         string                `gorm:"column:redirect_uri" json:"redirect_uri"`
	Scopes              string                `gorm:"column:scopes" json:"scopes"`
	CodeChallenge       *string               `gorm:"column:code_challenge" json:"code_challenge"`
	CodeChallengeMethod *string               `gorm:"column:code_challenge_method" json:"code_challenge_method"`
	ResourceOwner       Users                 `gorm:"foreignKey:ResourceOwnerId;references:ID"`
	Application         OauthApplications     `gorm:"foreignKey:ApplicationId;references:ID"`
	OauthOpenidRequests []OauthOpenidRequests `gorm:"foreignKey:AccessGrantId"`
}

// TableName sets the insert table name for this struct type
func (m *OauthAccessGrants) TableName() string {
	return "oauth_access_grants"
}
