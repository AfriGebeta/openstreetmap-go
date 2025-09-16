package models

type CurrentRelationMembers struct {
	RelationId int64            `gorm:"column:relation_id" json:"relation_id"`
	MemberType string           `gorm:"column:member_type" json:"member_type"`
	MemberId   int64            `gorm:"column:member_id" json:"member_id"`
	SequenceId int              `gorm:"primaryKey;column:sequence_id" json:"sequence_id"`
	MemberRole string           `gorm:"column:member_role" json:"member_role"`
	Relation   CurrentRelations `gorm:"foreignKey:RelationId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *CurrentRelationMembers) TableName() string {
	return "current_relation_members"
}
