package models

import "time"

type ChangesetComments struct {
	ID          int        `gorm:"primaryKey;column:id" json:"id"`
	ChangesetId int64      `gorm:"column:changeset_id" json:"changeset_id"`
	AuthorId    int64      `gorm:"column:author_id" json:"author_id"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	Visible     bool       `gorm:"column:visible" json:"visible"`
	Body        string     `gorm:"column:body" json:"body"`
	Author      Users      `gorm:"foreignKey:AuthorId;references:ID"`
	Changeset   Changesets `gorm:"foreignKey:ChangesetId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *ChangesetComments) TableName() string {
	return "changeset_comments"
}
