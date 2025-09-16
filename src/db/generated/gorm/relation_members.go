package models

type RelationMembers struct {
	RelationId int64     `gorm:"column:relation_id" json:"relation_id"`
	MemberType string    `gorm:"column:member_type" json:"member_type"`
	MemberId   int64     `gorm:"column:member_id" json:"member_id"`
	Version    int64     `gorm:"column:version" json:"version"`
	SequenceId int       `gorm:"primaryKey;column:sequence_id" json:"sequence_id"`
	MemberRole string    `gorm:"column:member_role" json:"member_role"`
	Relation   Relations `gorm:"foreignKey:RelationId;references:RelationId"`
}

// TableName sets the insert table name for this struct type
func (m *RelationMembers) TableName() string {
	return "relation_members"
}
