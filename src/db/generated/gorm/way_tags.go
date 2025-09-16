package models

type WayTags struct {
	Version int64  `gorm:"column:version" json:"version"`
	WayId   int64  `gorm:"column:way_id" json:"way_id"`
	K       string `gorm:"primaryKey;column:k" json:"k"`
	V       string `gorm:"column:v" json:"v"`
	Way     Ways   `gorm:"foreignKey:WayId;references:WayId"`
}

// TableName sets the insert table name for this struct type
func (m *WayTags) TableName() string {
	return "way_tags"
}
