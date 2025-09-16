package models

import "time"

type DiaryComments struct {
	ID           int64        `gorm:"primaryKey;column:id" json:"id"`
	DiaryEntryId int64        `gorm:"column:diary_entry_id" json:"diary_entry_id"`
	UserId       int64        `gorm:"column:user_id" json:"user_id"`
	CreatedAt    time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"column:updated_at" json:"updated_at"`
	Visible      bool         `gorm:"column:visible" json:"visible"`
	BodyFormat   string       `gorm:"column:body_format" json:"body_format"`
	Body         string       `gorm:"column:body" json:"body"`
	DiaryEntry   DiaryEntries `gorm:"foreignKey:DiaryEntryId;references:ID"`
	User         Users        `gorm:"foreignKey:UserId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *DiaryComments) TableName() string {
	return "diary_comments"
}
