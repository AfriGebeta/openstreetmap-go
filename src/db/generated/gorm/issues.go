package models

import "time"

type Issues struct {
	ID             int             `gorm:"primaryKey;column:id" json:"id"`
	ReportableId   int             `gorm:"column:reportable_id" json:"reportable_id"`
	ReportedUserId *int            `gorm:"column:reported_user_id" json:"reported_user_id"`
	Status         string          `gorm:"column:status" json:"status"`
	AssignedRole   string          `gorm:"column:assigned_role" json:"assigned_role"`
	ResolvedAt     *time.Time      `gorm:"column:resolved_at" json:"resolved_at"`
	ResolvedBy     *int            `gorm:"column:resolved_by" json:"resolved_by"`
	UpdatedBy      *int            `gorm:"column:updated_by" json:"updated_by"`
	ReportsCount   *int            `gorm:"column:reports_count" json:"reports_count"`
	CreatedAt      time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"column:updated_at" json:"updated_at"`
	ReportableType string          `gorm:"column:reportable_type" json:"reportable_type"`
	ReportedUser   Users           `gorm:"foreignKey:ReportedUserId;references:ID"`
	Resolved       Users           `gorm:"foreignKey:ResolvedBy;references:ID"`
	Updated        Users           `gorm:"foreignKey:UpdatedBy;references:ID"`
	IssueComments  []IssueComments `gorm:"foreignKey:IssueId"`
	Reports        []Reports       `gorm:"foreignKey:IssueId"`
}

// TableName sets the insert table name for this struct type
func (m *Issues) TableName() string {
	return "issues"
}
