package models

import "time"

type Redactions struct {
	ID                int         `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt         *time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         *time.Time  `gorm:"column:updated_at" json:"updated_at"`
	UserId            int64       `gorm:"column:user_id" json:"user_id"`
	DescriptionFormat string      `gorm:"column:description_format" json:"description_format"`
	Title             string      `gorm:"column:title" json:"title"`
	Description       string      `gorm:"column:description" json:"description"`
	User              Users       `gorm:"foreignKey:UserId;references:ID"`
	Nodes             []Nodes     `gorm:"foreignKey:RedactionId"`
	Relations         []Relations `gorm:"foreignKey:RedactionId"`
	Ways              []Ways      `gorm:"foreignKey:RedactionId"`
}

// TableName sets the insert table name for this struct type
func (m *Redactions) TableName() string {
	return "redactions"
}
