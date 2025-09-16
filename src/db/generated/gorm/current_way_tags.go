package models

type CurrentWayTags struct {
	WayId int64       `gorm:"column:way_id" json:"way_id"`
	K     string      `gorm:"primaryKey;column:k" json:"k"`
	V     string      `gorm:"column:v" json:"v"`
	Way   CurrentWays `gorm:"foreignKey:WayId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *CurrentWayTags) TableName() string {
	return "current_way_tags"
}
