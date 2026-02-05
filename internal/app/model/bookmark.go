package model

// Bookmark represents a bookmark entity in the system.
type Bookmark struct {
	Base

	Description string `gorm:"column:description" json:"description"`
	URL         string `gorm:"not null;column:url" json:"url"`
	Code        string `gorm:"unique;not null;column:code" json:"code"`
	UserID      string `gorm:"not null;column:user_id" json:"-"`
	User        User   `gorm:"foreignKey:UserID;references:ID" json:"-"`
}
