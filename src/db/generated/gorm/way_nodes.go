package models

type WayNodes struct {
	WayId      int64 `gorm:"column:way_id" json:"way_id"`
	NodeId     int64 `gorm:"column:node_id" json:"node_id"`
	Version    int64 `gorm:"column:version" json:"version"`
	SequenceId int64 `gorm:"primaryKey;column:sequence_id" json:"sequence_id"`
	Way        Ways  `gorm:"foreignKey:WayId;references:WayId"`
}

// TableName sets the insert table name for this struct type
func (m *WayNodes) TableName() string {
	return "way_nodes"
}
