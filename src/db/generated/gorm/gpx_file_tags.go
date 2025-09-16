package models

type GpxFileTags struct {
	GpxId int64    `gorm:"column:gpx_id" json:"gpx_id"`
	ID    int64    `gorm:"primaryKey;column:id" json:"id"`
	Tag   string   `gorm:"column:tag" json:"tag"`
	Gpx   GpxFiles `gorm:"foreignKey:GpxId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *GpxFileTags) TableName() string {
	return "gpx_file_tags"
}
