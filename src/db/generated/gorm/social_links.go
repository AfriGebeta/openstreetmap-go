package models

import "time"

type SocialLinks struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"id"`
	UserId    int64     `gorm:"column:user_id" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Url       string    `gorm:"column:url" json:"url"`
	User      Users     `gorm:"foreignKey:UserId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *SocialLinks) TableName() string {
	return "social_links"
}
