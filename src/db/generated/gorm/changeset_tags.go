package models

type ChangesetTags struct {
	ChangesetId int64      `gorm:"column:changeset_id" json:"changeset_id"`
	K           string     `gorm:"primaryKey;column:k" json:"k"`
	V           string     `gorm:"column:v" json:"v"`
	Changeset   Changesets `gorm:"foreignKey:ChangesetId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *ChangesetTags) TableName() string {
	return "changeset_tags"
}
