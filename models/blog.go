package models

import (
	"time"
)

type Blog struct {
	ID       uint      `gorm:"primaryKey"`
	AuthorID uint      `gorm:"column:author_id"`
	Title    string    `gorm:"column:title"`
	Content  string    `gorm:"column:content"`
	Markdown bool      `gorm:"column:markdown"`
	PostedOn time.Time `gorm:"column:posted_on"`

	Tags      []Tag      `gorm:"many2many:has_tags"`
	Reactions []Reaction `gorm:"foreignKey:blog_id"`
	Comments  []Comment  `gorm:"foreignKey:blog_id"`
}

func (*Blog) TableName() string {
	return "blog"
}
