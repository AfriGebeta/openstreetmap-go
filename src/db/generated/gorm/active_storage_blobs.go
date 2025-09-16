package models

import "time"

type ActiveStorageBlobs struct {
	ID                          int64                         `gorm:"primaryKey;column:id" json:"id"`
	ByteSize                    int64                         `gorm:"column:byte_size" json:"byte_size"`
	CreatedAt                   time.Time                     `gorm:"column:created_at" json:"created_at"`
	Key                         string                        `gorm:"column:key" json:"key"`
	Filename                    string                        `gorm:"column:filename" json:"filename"`
	ContentType                 *string                       `gorm:"column:content_type" json:"content_type"`
	Metadata                    *string                       `gorm:"column:metadata" json:"metadata"`
	Checksum                    *string                       `gorm:"column:checksum" json:"checksum"`
	ServiceName                 string                        `gorm:"column:service_name" json:"service_name"`
	ActiveStorageAttachments    []ActiveStorageAttachments    `gorm:"foreignKey:BlobId"`
	ActiveStorageVariantRecords []ActiveStorageVariantRecords `gorm:"foreignKey:BlobId"`
}

// TableName sets the insert table name for this struct type
func (m *ActiveStorageBlobs) TableName() string {
	return "active_storage_blobs"
}
