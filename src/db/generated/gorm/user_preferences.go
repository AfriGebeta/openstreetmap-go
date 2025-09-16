package models

type UserPreferences struct {
	UserId int64  `gorm:"column:user_id" json:"user_id"`
	K      string `gorm:"primaryKey;column:k" json:"k"`
	V      string `gorm:"column:v" json:"v"`
	User   Users  `gorm:"foreignKey:UserId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *UserPreferences) TableName() string {
	return "user_preferences"
}
