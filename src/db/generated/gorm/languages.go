package models

type Languages struct {
	Code         string         `gorm:"primaryKey;column:code" json:"code"`
	EnglishName  string         `gorm:"column:english_name" json:"english_name"`
	NativeName   *string        `gorm:"column:native_name" json:"native_name"`
	DiaryEntries []DiaryEntries `gorm:"foreignKey:LanguageCode"`
}

// TableName sets the insert table name for this struct type
func (m *Languages) TableName() string {
	return "languages"
}
