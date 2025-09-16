package models

import "time"

type Nodes struct {
	NodeId      int64      `gorm:"column:node_id" json:"node_id"`
	Latitude    int        `gorm:"column:latitude" json:"latitude"`
	Longitude   int        `gorm:"column:longitude" json:"longitude"`
	ChangesetId int64      `gorm:"column:changeset_id" json:"changeset_id"`
	Visible     bool       `gorm:"column:visible" json:"visible"`
	Timestamp   time.Time  `gorm:"column:timestamp" json:"timestamp"`
	Tile        int64      `gorm:"column:tile" json:"tile"`
	Version     int64      `gorm:"primaryKey;column:version" json:"version"`
	RedactionId *int       `gorm:"column:redaction_id" json:"redaction_id"`
	Changeset   Changesets `gorm:"foreignKey:ChangesetId;references:ID"`
	Redaction   Redactions `gorm:"foreignKey:RedactionId;references:ID"`
	NodeTags    []NodeTags `gorm:"foreignKey:NodeId"`
}

// TableName sets the insert table name for this struct type
func (m *Nodes) TableName() string {
	return "nodes"
}
