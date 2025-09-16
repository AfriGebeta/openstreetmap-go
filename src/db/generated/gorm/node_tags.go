package models

type NodeTags struct {
	NodeId  int64  `gorm:"column:node_id" json:"node_id"`
	Version int64  `gorm:"column:version" json:"version"`
	K       string `gorm:"primaryKey;column:k" json:"k"`
	V       string `gorm:"column:v" json:"v"`
	Node    Nodes  `gorm:"foreignKey:NodeId;references:NodeId"`
}

// TableName sets the insert table name for this struct type
func (m *NodeTags) TableName() string {
	return "node_tags"
}
