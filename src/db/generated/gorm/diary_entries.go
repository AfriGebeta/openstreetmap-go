package models

import "time"

type DiaryEntries struct {
	ID                      int64                     `gorm:"primaryKey;column:id" json:"id"`
	UserId                  int64                     `gorm:"column:user_id" json:"user_id"`
	CreatedAt               time.Time                 `gorm:"column:created_at" json:"created_at"`
	UpdatedAt               time.Time                 `gorm:"column:updated_at" json:"updated_at"`
	Latitude                *float64                  `gorm:"column:latitude" json:"latitude"`
	Longitude               *float64                  `gorm:"column:longitude" json:"longitude"`
	Visible                 bool                      `gorm:"column:visible" json:"visible"`
	BodyFormat              string                    `gorm:"column:body_format" json:"body_format"`
	Title                   string                    `gorm:"column:title" json:"title"`
	Body                    string                    `gorm:"column:body" json:"body"`
	LanguageCode            string                    `gorm:"column:language_code" json:"language_code"`
	Language                Languages                 `gorm:"foreignKey:LanguageCode;references:Code"`
	User                    Users                     `gorm:"foreignKey:UserId;references:ID"`
	DiaryComments           []DiaryComments           `gorm:"foreignKey:DiaryEntryId"`
	DiaryEntrySubscriptions []DiaryEntrySubscriptions `gorm:"foreignKey:DiaryEntryId"`
}

// TableName sets the insert table name for this struct type
func (m *DiaryEntries) TableName() string {
	return "diary_entries"
}
