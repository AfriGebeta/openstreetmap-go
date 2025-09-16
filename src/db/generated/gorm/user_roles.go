package models

import "time"

type UserRoles struct {
	ID        int        `gorm:"primaryKey;column:id" json:"id"`
	UserId    int64      `gorm:"column:user_id" json:"user_id"`
	Role      string     `gorm:"column:role" json:"role"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	GranterId int64      `gorm:"column:granter_id" json:"granter_id"`
	Granter   Users      `gorm:"foreignKey:GranterId;references:ID"`
	User      Users      `gorm:"foreignKey:UserId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *UserRoles) TableName() string {
	return "user_roles"
}
