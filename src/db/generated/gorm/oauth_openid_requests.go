package models

type OauthOpenidRequests struct {
	ID            int64             `gorm:"primaryKey;column:id" json:"id"`
	AccessGrantId int64             `gorm:"column:access_grant_id" json:"access_grant_id"`
	Nonce         string            `gorm:"column:nonce" json:"nonce"`
	AccessGrant   OauthAccessGrants `gorm:"foreignKey:AccessGrantId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *OauthOpenidRequests) TableName() string {
	return "oauth_openid_requests"
}
