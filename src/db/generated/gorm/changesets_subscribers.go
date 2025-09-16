package models

type ChangesetsSubscribers struct {
	SubscriberId int64      `gorm:"column:subscriber_id" json:"subscriber_id"`
	ChangesetId  int64      `gorm:"column:changeset_id" json:"changeset_id"`
	Changeset    Changesets `gorm:"foreignKey:ChangesetId;references:ID"`
	Subscriber   Users      `gorm:"foreignKey:SubscriberId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *ChangesetsSubscribers) TableName() string {
	return "changesets_subscribers"
}
