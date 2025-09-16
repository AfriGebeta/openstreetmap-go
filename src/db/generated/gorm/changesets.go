package models

import "time"

type Changesets struct {
	ID                    int64                   `gorm:"primaryKey;column:id" json:"id"`
	UserId                int64                   `gorm:"column:user_id" json:"user_id"`
	CreatedAt             time.Time               `gorm:"column:created_at" json:"created_at"`
	MinLat                *int                    `gorm:"column:min_lat" json:"min_lat"`
	MaxLat                *int                    `gorm:"column:max_lat" json:"max_lat"`
	MinLon                *int                    `gorm:"column:min_lon" json:"min_lon"`
	MaxLon                *int                    `gorm:"column:max_lon" json:"max_lon"`
	ClosedAt              time.Time               `gorm:"column:closed_at" json:"closed_at"`
	NumChanges            int                     `gorm:"column:num_changes" json:"num_changes"`
	NumCreatedNodes       int                     `gorm:"column:num_created_nodes" json:"num_created_nodes"`
	NumModifiedNodes      int                     `gorm:"column:num_modified_nodes" json:"num_modified_nodes"`
	NumDeletedNodes       int                     `gorm:"column:num_deleted_nodes" json:"num_deleted_nodes"`
	NumCreatedWays        int                     `gorm:"column:num_created_ways" json:"num_created_ways"`
	NumModifiedWays       int                     `gorm:"column:num_modified_ways" json:"num_modified_ways"`
	NumDeletedWays        int                     `gorm:"column:num_deleted_ways" json:"num_deleted_ways"`
	NumCreatedRelations   int                     `gorm:"column:num_created_relations" json:"num_created_relations"`
	NumModifiedRelations  int                     `gorm:"column:num_modified_relations" json:"num_modified_relations"`
	NumDeletedRelations   int                     `gorm:"column:num_deleted_relations" json:"num_deleted_relations"`
	User                  Users                   `gorm:"foreignKey:UserId;references:ID"`
	ChangesetComments     []ChangesetComments     `gorm:"foreignKey:ChangesetId"`
	ChangesetTags         []ChangesetTags         `gorm:"foreignKey:ChangesetId"`
	ChangesetsSubscribers []ChangesetsSubscribers `gorm:"foreignKey:ChangesetId"`
	CurrentNodes          []CurrentNodes          `gorm:"foreignKey:ChangesetId"`
	CurrentRelations      []CurrentRelations      `gorm:"foreignKey:ChangesetId"`
	CurrentWays           []CurrentWays           `gorm:"foreignKey:ChangesetId"`
	Nodes                 []Nodes                 `gorm:"foreignKey:ChangesetId"`
	Relations             []Relations             `gorm:"foreignKey:ChangesetId"`
	Ways                  []Ways                  `gorm:"foreignKey:ChangesetId"`
}

// TableName sets the insert table name for this struct type
func (m *Changesets) TableName() string {
	return "changesets"
}
