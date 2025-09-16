package models

import "time"

type GpxFiles struct {
	ID          int64         `gorm:"primaryKey;column:id" json:"id"`
	UserId      int64         `gorm:"column:user_id" json:"user_id"`
	Visible     bool          `gorm:"column:visible" json:"visible"`
	Size        *int64        `gorm:"column:size" json:"size"`
	Latitude    *float64      `gorm:"column:latitude" json:"latitude"`
	Longitude   *float64      `gorm:"column:longitude" json:"longitude"`
	Timestamp   time.Time     `gorm:"column:timestamp" json:"timestamp"`
	Inserted    bool          `gorm:"column:inserted" json:"inserted"`
	Visibility  string        `gorm:"column:visibility" json:"visibility"`
	Name        string        `gorm:"column:name" json:"name"`
	Description string        `gorm:"column:description" json:"description"`
	User        Users         `gorm:"foreignKey:UserId;references:ID"`
	GpsPoints   []GpsPoints   `gorm:"foreignKey:GpxId"`
	GpxFileTags []GpxFileTags `gorm:"foreignKey:GpxId"`
}

// TableName sets the insert table name for this struct type
func (m *GpxFiles) TableName() string {
	return "gpx_files"
}
