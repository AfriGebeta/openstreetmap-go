package models

import "time"

type Ways struct {
	WayId       int64      `gorm:"column:way_id" json:"way_id"`
	ChangesetId int64      `gorm:"column:changeset_id" json:"changeset_id"`
	Timestamp   time.Time  `gorm:"column:timestamp" json:"timestamp"`
	Version     int64      `gorm:"primaryKey;column:version" json:"version"`
	Visible     bool       `gorm:"column:visible" json:"visible"`
	RedactionId *int       `gorm:"column:redaction_id" json:"redaction_id"`
	Changeset   Changesets `gorm:"foreignKey:ChangesetId;references:ID"`
	Redaction   Redactions `gorm:"foreignKey:RedactionId;references:ID"`
	WayNodes    []WayNodes `gorm:"foreignKey:WayId"`
	WayTags     []WayTags  `gorm:"foreignKey:WayId"`
}

// TableName sets the insert table name for this struct type
func (m *Ways) TableName() string {
	return "ways"
}
