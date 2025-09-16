package models

type RelationTags struct {
	RelationId int64     `gorm:"column:relation_id" json:"relation_id"`
	Version    int64     `gorm:"column:version" json:"version"`
	V          string    `gorm:"column:v" json:"v"`
	K          string    `gorm:"primaryKey;column:k" json:"k"`
	Relation   Relations `gorm:"foreignKey:RelationId;references:RelationId"`
}

// TableName sets the insert table name for this struct type
func (m *RelationTags) TableName() string {
	return "relation_tags"
}
