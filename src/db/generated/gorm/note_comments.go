package models

import "time"

type NoteComments struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"id"`
	NoteId    int64     `gorm:"column:note_id" json:"note_id"`
	Visible   bool      `gorm:"column:visible" json:"visible"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	AuthorIp  *string   `gorm:"column:author_ip" json:"author_ip"`
	AuthorId  *int64    `gorm:"column:author_id" json:"author_id"`
	Event     *string   `gorm:"column:event" json:"event"`
	Body      *string   `gorm:"column:body" json:"body"`
	Author    Users     `gorm:"foreignKey:AuthorId;references:ID"`
	Note      Notes     `gorm:"foreignKey:NoteId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *NoteComments) TableName() string {
	return "note_comments"
}
