package models

import "time"

type Relations struct {
	RelationId      int64             `gorm:"column:relation_id" json:"relation_id"`
	ChangesetId     int64             `gorm:"column:changeset_id" json:"changeset_id"`
	Timestamp       time.Time         `gorm:"column:timestamp" json:"timestamp"`
	Version         int64             `gorm:"primaryKey;column:version" json:"version"`
	Visible         bool              `gorm:"column:visible" json:"visible"`
	RedactionId     *int              `gorm:"column:redaction_id" json:"redaction_id"`
	Changeset       Changesets        `gorm:"foreignKey:ChangesetId;references:ID"`
	Redaction       Redactions        `gorm:"foreignKey:RedactionId;references:ID"`
	RelationMembers []RelationMembers `gorm:"foreignKey:RelationId"`
	RelationTags    []RelationTags    `gorm:"foreignKey:RelationId"`
}

// TableName sets the insert table name for this struct type
func (m *Relations) TableName() string {
	return "relations"
}
