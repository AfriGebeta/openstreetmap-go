package models

type CurrentWayNodes struct {
	WayId      int64        `gorm:"column:way_id" json:"way_id"`
	NodeId     int64        `gorm:"column:node_id" json:"node_id"`
	SequenceId int64        `gorm:"primaryKey;column:sequence_id" json:"sequence_id"`
	Way        CurrentWays  `gorm:"foreignKey:WayId;references:ID"`
	Node       CurrentNodes `gorm:"foreignKey:NodeId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *CurrentWayNodes) TableName() string {
	return "current_way_nodes"
}
