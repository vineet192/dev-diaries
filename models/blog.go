package models

import (
	"time"
)

type Blog struct {
	ID       uint      `gorm:"primaryKey;column:id"`
	AuthorID uint      `gorm:"column:author_id"`
	Title    string    `gorm:"column:title"`
	Content  string    `gorm:"column:content"`
	Markdown bool      `gorm:"column:markdown"`
	PostedOn time.Time `gorm:"column:posted_on"`

	Tags      []Tag      `gorm:"many2many:has_tags;foreignKey:blog_id"`
	Reactions []Reaction `gorm:"foreignKey:blog_id"`
	Comments  []Comment  `gorm:"foreignKey:blog_id"`
}
