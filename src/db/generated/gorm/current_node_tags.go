package models

type CurrentNodeTags struct {
	NodeId int64        `gorm:"column:node_id" json:"node_id"`
	K      string       `gorm:"primaryKey;column:k" json:"k"`
	V      string       `gorm:"column:v" json:"v"`
	Node   CurrentNodes `gorm:"foreignKey:NodeId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *CurrentNodeTags) TableName() string {
	return "current_node_tags"
}
