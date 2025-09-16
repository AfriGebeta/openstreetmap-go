package models

import "time"

type UserBlocks struct {
	ID            int        `gorm:"primaryKey;column:id" json:"id"`
	UserId        int64      `gorm:"column:user_id" json:"user_id"`
	CreatorId     int64      `gorm:"column:creator_id" json:"creator_id"`
	EndsAt        time.Time  `gorm:"column:ends_at" json:"ends_at"`
	NeedsView     bool       `gorm:"column:needs_view" json:"needs_view"`
	RevokerId     *int64     `gorm:"column:revoker_id" json:"revoker_id"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`
	ReasonFormat  string     `gorm:"column:reason_format" json:"reason_format"`
	DeactivatesAt *time.Time `gorm:"column:deactivates_at" json:"deactivates_at"`
	Reason        string     `gorm:"column:reason" json:"reason"`
	Creator       Users      `gorm:"foreignKey:CreatorId;references:ID"`
	Revoker       Users      `gorm:"foreignKey:RevokerId;references:ID"`
	User          Users      `gorm:"foreignKey:UserId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *UserBlocks) TableName() string {
	return "user_blocks"
}
