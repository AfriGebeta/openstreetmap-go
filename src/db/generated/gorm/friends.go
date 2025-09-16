package models

import "time"

type Friends struct {
	ID           int64      `gorm:"primaryKey;column:id" json:"id"`
	UserId       int64      `gorm:"column:user_id" json:"user_id"`
	FriendUserId int64      `gorm:"column:friend_user_id" json:"friend_user_id"`
	CreatedAt    *time.Time `gorm:"column:created_at" json:"created_at"`
	FriendUser   Users      `gorm:"foreignKey:FriendUserId;references:ID"`
	User         Users      `gorm:"foreignKey:UserId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *Friends) TableName() string {
	return "friends"
}
