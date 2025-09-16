package models

import "time"

type Messages struct {
	ID              int64     `gorm:"primaryKey;column:id" json:"id"`
	FromUserId      int64     `gorm:"column:from_user_id" json:"from_user_id"`
	SentOn          time.Time `gorm:"column:sent_on" json:"sent_on"`
	MessageRead     bool      `gorm:"column:message_read" json:"message_read"`
	ToUserId        int64     `gorm:"column:to_user_id" json:"to_user_id"`
	ToUserVisible   bool      `gorm:"column:to_user_visible" json:"to_user_visible"`
	FromUserVisible bool      `gorm:"column:from_user_visible" json:"from_user_visible"`
	BodyFormat      string    `gorm:"column:body_format" json:"body_format"`
	Muted           bool      `gorm:"column:muted" json:"muted"`
	Title           string    `gorm:"column:title" json:"title"`
	Body            string    `gorm:"column:body" json:"body"`
	FromUser        Users     `gorm:"foreignKey:FromUserId;references:ID"`
	ToUser          Users     `gorm:"foreignKey:ToUserId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *Messages) TableName() string {
	return "messages"
}
