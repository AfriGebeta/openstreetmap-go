package models

import "time"

type UserMutes struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"id"`
	OwnerId   int64     `gorm:"column:owner_id" json:"owner_id"`
	SubjectId int64     `gorm:"column:subject_id" json:"subject_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Owner     Users     `gorm:"foreignKey:OwnerId;references:ID"`
	Subject   Users     `gorm:"foreignKey:SubjectId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *UserMutes) TableName() string {
	return "user_mutes"
}
