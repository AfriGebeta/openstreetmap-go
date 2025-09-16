package models

import "time"

type DelayedJobs struct {
	ID        int64      `gorm:"primaryKey;column:id" json:"id"`
	Priority  int        `gorm:"column:priority" json:"priority"`
	Attempts  int        `gorm:"column:attempts" json:"attempts"`
	RunAt     *time.Time `gorm:"column:run_at" json:"run_at"`
	LockedAt  *time.Time `gorm:"column:locked_at" json:"locked_at"`
	FailedAt  *time.Time `gorm:"column:failed_at" json:"failed_at"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	Handler   string     `gorm:"column:handler" json:"handler"`
	LastError *string    `gorm:"column:last_error" json:"last_error"`
	LockedBy  *string    `gorm:"column:locked_by" json:"locked_by"`
	Queue     *string    `gorm:"column:queue" json:"queue"`
}

// TableName sets the insert table name for this struct type
func (m *DelayedJobs) TableName() string {
	return "delayed_jobs"
}
