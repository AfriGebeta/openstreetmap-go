package models

type SchemaMigrations struct {
	Version string `gorm:"primaryKey;column:version" json:"version"`
}

// TableName sets the insert table name for this struct type
func (m *SchemaMigrations) TableName() string {
	return "schema_migrations"
}
