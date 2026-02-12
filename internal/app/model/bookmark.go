package model

// Bookmark represents a bookmark entity in the system.
type Bookmark struct {
	Base

	Description string `gorm:"column:description" json:"description"`
	URL         string `gorm:"not null;column:url" json:"url"`
	Code        string `gorm:"column:code" json:"-"` // this field will be deleted form version 3 and replaced by CodeShortenEncoded

	UserID string `gorm:"not null;column:user_id" json:"-"`
	User   User   `gorm:"foreignKey:UserID;references:ID" json:"-"`

	CodeShorten        int64  `gorm:"column:code_shorten;autoIncrement" json:"-"`
	CodeShortenEncoded string `gorm:"column:code_shorten_encoded" json:"code"`
}
