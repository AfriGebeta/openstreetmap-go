package models

import "time"

type ArInternalMetadata struct {
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Key       string    `gorm:"primaryKey;column:key" json:"key"`
	Value     *string   `gorm:"column:value" json:"value"`
}

// TableName sets the insert table name for this struct type
func (m *ArInternalMetadata) TableName() string {
	return "ar_internal_metadata"
}
