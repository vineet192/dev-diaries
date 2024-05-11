package models

import (
	"time"
)

type Blog struct {
	id        uint `gorm:"primaryKey"`
	author_id uint
	title     string
	content   string
	markdown  string
	posted_on time.Time

	Tags      []Tag      `gorm:"many2many:has_tags;foreignKey:blog_id"`
	Reactions []Reaction `gorm:"foreignKey:blog_id"`
	Comments  []Comment  `gorm:"foreignKey:blog_id"`
}
