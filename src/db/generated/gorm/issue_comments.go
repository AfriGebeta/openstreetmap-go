package models

import "time"

type IssueComments struct {
	ID        int       `gorm:"primaryKey;column:id" json:"id"`
	IssueId   int       `gorm:"column:issue_id" json:"issue_id"`
	UserId    int       `gorm:"column:user_id" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Body      string    `gorm:"column:body" json:"body"`
	Issue     Issues    `gorm:"foreignKey:IssueId;references:ID"`
	User      Users     `gorm:"foreignKey:UserId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *IssueComments) TableName() string {
	return "issue_comments"
}
