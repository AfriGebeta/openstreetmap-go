package models

type Acls struct {
	Address *string `gorm:"column:address" json:"address"`
	ID      int64   `gorm:"primaryKey;column:id" json:"id"`
	K       string  `gorm:"column:k" json:"k"`
	V       *string `gorm:"column:v" json:"v"`
	Domain  *string `gorm:"column:domain" json:"domain"`
	Mx      *string `gorm:"column:mx" json:"mx"`
}

// TableName sets the insert table name for this struct type
func (m *Acls) TableName() string {
	return "acls"
}
