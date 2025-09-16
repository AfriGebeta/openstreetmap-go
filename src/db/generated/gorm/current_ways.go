package models

import "time"

type CurrentWays struct {
	ID              int64             `gorm:"primaryKey;column:id" json:"id"`
	ChangesetId     int64             `gorm:"column:changeset_id" json:"changeset_id"`
	Timestamp       time.Time         `gorm:"column:timestamp" json:"timestamp"`
	Visible         bool              `gorm:"column:visible" json:"visible"`
	Version         int64             `gorm:"column:version" json:"version"`
	Changeset       Changesets        `gorm:"foreignKey:ChangesetId;references:ID"`
	CurrentWayNodes []CurrentWayNodes `gorm:"foreignKey:WayId"`
	CurrentWayTags  []CurrentWayTags  `gorm:"foreignKey:WayId"`
}

// TableName sets the insert table name for this struct type
func (m *CurrentWays) TableName() string {
	return "current_ways"
}
