package models

type NoteSubscriptions struct {
	UserId int64 `gorm:"column:user_id" json:"user_id"`
	NoteId int64 `gorm:"primaryKey;column:note_id" json:"note_id"`
	Note   Notes `gorm:"foreignKey:NoteId;references:ID"`
	User   Users `gorm:"foreignKey:UserId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *NoteSubscriptions) TableName() string {
	return "note_subscriptions"
}
