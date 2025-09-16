package models

import "time"

type GpsPoints struct {
	Altitude  *float64   `gorm:"column:altitude" json:"altitude"`
	Trackid   int        `gorm:"column:trackid" json:"trackid"`
	Latitude  int        `gorm:"column:latitude" json:"latitude"`
	Longitude int        `gorm:"column:longitude" json:"longitude"`
	GpxId     int64      `gorm:"column:gpx_id" json:"gpx_id"`
	Timestamp *time.Time `gorm:"column:timestamp" json:"timestamp"`
	Tile      *int64     `gorm:"column:tile" json:"tile"`
	Gpx       GpxFiles   `gorm:"foreignKey:GpxId;references:ID"`
}

// TableName sets the insert table name for this struct type
func (m *GpsPoints) TableName() string {
	return "gps_points"
}
