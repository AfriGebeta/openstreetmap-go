package models

import "time"

type CurrentNodes struct {
	ID              int64             `gorm:"primaryKey;column:id" json:"id"`
	Latitude        int               `gorm:"column:latitude" json:"latitude"`
	Longitude       int               `gorm:"column:longitude" json:"longitude"`
	ChangesetId     int64             `gorm:"column:changeset_id" json:"changeset_id"`
	Visible         bool              `gorm:"column:visible" json:"visible"`
	Timestamp       time.Time         `gorm:"column:timestamp" json:"timestamp"`
	Tile            int64             `gorm:"column:tile" json:"tile"`
	Version         int64             `gorm:"column:version" json:"version"`
	Changeset       Changesets        `gorm:"foreignKey:ChangesetId;references:ID"`
	CurrentNodeTags []CurrentNodeTags `gorm:"foreignKey:NodeId"`
	CurrentWayNodes []CurrentWayNodes `gorm:"foreignKey:NodeId"`
}

// TableName sets the insert table name for this struct type
func (m *CurrentNodes) TableName() string {
	return "current_nodes"
}
