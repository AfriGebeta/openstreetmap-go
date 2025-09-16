package models

import "time"

type Users struct {
	ID                      int64                     `gorm:"primaryKey;column:id" json:"id"`
	CreationTime            time.Time                 `gorm:"column:creation_time" json:"creation_time"`
	DataPublic              bool                      `gorm:"column:data_public" json:"data_public"`
	HomeLat                 *float64                  `gorm:"column:home_lat" json:"home_lat"`
	HomeLon                 *float64                  `gorm:"column:home_lon" json:"home_lon"`
	HomeZoom                *int                      `gorm:"column:home_zoom" json:"home_zoom"`
	EmailValid              bool                      `gorm:"column:email_valid" json:"email_valid"`
	Status                  string                    `gorm:"column:status" json:"status"`
	TermsAgreed             *time.Time                `gorm:"column:terms_agreed" json:"terms_agreed"`
	ConsiderPd              bool                      `gorm:"column:consider_pd" json:"consider_pd"`
	TermsSeen               bool                      `gorm:"column:terms_seen" json:"terms_seen"`
	DescriptionFormat       string                    `gorm:"column:description_format" json:"description_format"`
	ChangesetsCount         int                       `gorm:"column:changesets_count" json:"changesets_count"`
	TracesCount             int                       `gorm:"column:traces_count" json:"traces_count"`
	DiaryEntriesCount       int                       `gorm:"column:diary_entries_count" json:"diary_entries_count"`
	ImageUseGravatar        bool                      `gorm:"column:image_use_gravatar" json:"image_use_gravatar"`
	HomeTile                *int64                    `gorm:"column:home_tile" json:"home_tile"`
	TouAgreed               *time.Time                `gorm:"column:tou_agreed" json:"tou_agreed"`
	DiaryCommentsCount      *int                      `gorm:"column:diary_comments_count" json:"diary_comments_count"`
	NoteCommentsCount       *int                      `gorm:"column:note_comments_count" json:"note_comments_count"`
	CreationAddress         *string                   `gorm:"column:creation_address" json:"creation_address"`
	HomeLocationName        *string                   `gorm:"column:home_location_name" json:"home_location_name"`
	Company                 *string                   `gorm:"column:company" json:"company"`
	Email                   string                    `gorm:"column:email" json:"email"`
	PassCrypt               string                    `gorm:"column:pass_crypt" json:"pass_crypt"`
	DisplayName             string                    `gorm:"column:display_name" json:"display_name"`
	Description             string                    `gorm:"column:description" json:"description"`
	PassSalt                *string                   `gorm:"column:pass_salt" json:"pass_salt"`
	NewEmail                *string                   `gorm:"column:new_email" json:"new_email"`
	Languages               *string                   `gorm:"column:languages" json:"languages"`
	AuthUid                 *string                   `gorm:"column:auth_uid" json:"auth_uid"`
	PreferredEditor         *string                   `gorm:"column:preferred_editor" json:"preferred_editor"`
	AuthProvider            *string                   `gorm:"column:auth_provider" json:"auth_provider"`
	ChangesetComments       []ChangesetComments       `gorm:"foreignKey:AuthorId"`
	Changesets              []Changesets              `gorm:"foreignKey:UserId"`
	ChangesetsSubscribers   []ChangesetsSubscribers   `gorm:"foreignKey:SubscriberId"`
	DiaryComments           []DiaryComments           `gorm:"foreignKey:UserId"`
	DiaryEntries            []DiaryEntries            `gorm:"foreignKey:UserId"`
	DiaryEntrySubscriptions []DiaryEntrySubscriptions `gorm:"foreignKey:UserId"`
	FriendsAsFriendUser     []Friends                 `gorm:"foreignKey:FriendUserId"`
	FriendsAsUser           []Friends                 `gorm:"foreignKey:UserId"`
	GpxFiles                []GpxFiles                `gorm:"foreignKey:UserId"`
	IssueComments           []IssueComments           `gorm:"foreignKey:UserId"`
	IssuesAsReportedUser    []Issues                  `gorm:"foreignKey:ReportedUserId"`
	IssuesAsResolvedBy      []Issues                  `gorm:"foreignKey:ResolvedBy"`
	IssuesAsUpdatedBy       []Issues                  `gorm:"foreignKey:UpdatedBy"`
	MessagesAsFromUser      []Messages                `gorm:"foreignKey:FromUserId"`
	MessagesAsToUser        []Messages                `gorm:"foreignKey:ToUserId"`
	NoteComments            []NoteComments            `gorm:"foreignKey:AuthorId"`
	NoteSubscriptions       []NoteSubscriptions       `gorm:"foreignKey:UserId"`
	Notes                   []Notes                   `gorm:"foreignKey:UserId"`
	OauthAccessGrants       []OauthAccessGrants       `gorm:"foreignKey:ResourceOwnerId"`
	OauthAccessTokens       []OauthAccessTokens       `gorm:"foreignKey:ResourceOwnerId"`
	OauthApplications       []OauthApplications       `gorm:"foreignKey:OwnerId"`
	Redactions              []Redactions              `gorm:"foreignKey:UserId"`
	Reports                 []Reports                 `gorm:"foreignKey:UserId"`
	SocialLinks             []SocialLinks             `gorm:"foreignKey:UserId"`
	UserBlocksAsCreator     []UserBlocks              `gorm:"foreignKey:CreatorId"`
	UserBlocksAsRevoker     []UserBlocks              `gorm:"foreignKey:RevokerId"`
	UserBlocksAsUser        []UserBlocks              `gorm:"foreignKey:UserId"`
	UserMutesAsOwner        []UserMutes               `gorm:"foreignKey:OwnerId"`
	UserMutesAsSubject      []UserMutes               `gorm:"foreignKey:SubjectId"`
	UserPreferences         []UserPreferences         `gorm:"foreignKey:UserId"`
	UserRolesAsGranter      []UserRoles               `gorm:"foreignKey:GranterId"`
	UserRolesAsUser         []UserRoles               `gorm:"foreignKey:UserId"`
}

// TableName sets the insert table name for this struct type
func (m *Users) TableName() string {
	return "users"
}
