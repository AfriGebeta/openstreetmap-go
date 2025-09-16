package models

import "time"

type Notes struct {
	ID                int64               `gorm:"primaryKey;column:id" json:"id"`
	Latitude          int                 `gorm:"column:latitude" json:"latitude"`
	Longitude         int                 `gorm:"column:longitude" json:"longitude"`
	Tile              int64               `gorm:"column:tile" json:"tile"`
	UpdatedAt         time.Time           `gorm:"column:updated_at" json:"updated_at"`
	CreatedAt         time.Time           `gorm:"column:created_at" json:"created_at"`
	Status            string              `gorm:"column:status" json:"status"`
	ClosedAt          *time.Time          `gorm:"column:closed_at" json:"closed_at"`
	UserId            *int64              `gorm:"column:user_id" json:"user_id"`
	UserIp            *string             `gorm:"column:user_ip" json:"user_ip"`
	Description       string              `gorm:"column:description" json:"description"`
	User              Users               `gorm:"foreignKey:UserId;references:ID"`
	NoteComments      []NoteComments      `gorm:"foreignKey:NoteId"`
	NoteSubscriptions []NoteSubscriptions `gorm:"foreignKey:NoteId"`
}

// TableName sets the insert table name for this struct type
func (m *Notes) TableName() string {
	return "notes"
}
