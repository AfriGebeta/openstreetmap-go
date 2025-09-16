package models

type CurrentRelationTags struct {
	RelationId int64            `gorm:"column:relation_id" json:"relation_id"`
	K          string           `gorm:"primaryKey;column:k" json:"k"`
	V          string           `gorm:"column:v" json:"v"`
	Relation   CurrentRelations `gorm:"foreignKey:RelationId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *CurrentRelationTags) TableName() string {
	return "current_relation_tags"
}
