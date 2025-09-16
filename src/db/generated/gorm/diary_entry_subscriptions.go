package models

type DiaryEntrySubscriptions struct {
	UserId       int64        `gorm:"column:user_id" json:"user_id"`
	DiaryEntryId int64        `gorm:"primaryKey;column:diary_entry_id" json:"diary_entry_id"`
	DiaryEntry   DiaryEntries `gorm:"foreignKey:DiaryEntryId;references:ID"`
	User         Users        `gorm:"foreignKey:UserId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *DiaryEntrySubscriptions) TableName() string {
	return "diary_entry_subscriptions"
}
