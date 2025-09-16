package models

type ActiveStorageVariantRecords struct {
	ID              int64              `gorm:"primaryKey;column:id" json:"id"`
	BlobId          int64              `gorm:"column:blob_id" json:"blob_id"`
	VariationDigest string             `gorm:"column:variation_digest" json:"variation_digest"`
	Blob            ActiveStorageBlobs `gorm:"foreignKey:BlobId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *ActiveStorageVariantRecords) TableName() string {
	return "active_storage_variant_records"
}
