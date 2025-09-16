package models

import "time"

type ActiveStorageAttachments struct {
	ID         int64              `gorm:"primaryKey;column:id" json:"id"`
	RecordId   int64              `gorm:"column:record_id" json:"record_id"`
	BlobId     int64              `gorm:"column:blob_id" json:"blob_id"`
	CreatedAt  time.Time          `gorm:"column:created_at" json:"created_at"`
	Name       string             `gorm:"column:name" json:"name"`
	RecordType string             `gorm:"column:record_type" json:"record_type"`
	Blob       ActiveStorageBlobs `gorm:"foreignKey:BlobId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *ActiveStorageAttachments) TableName() string {
	return "active_storage_attachments"
}
