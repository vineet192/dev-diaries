package models

import (
	"time"
)

type reaction_type string

type Reaction struct {
	ID           uint          `gorm:"primaryKey"`
	CommentID    uint          `gorm:"column:comment_id"`
	BlogID       uint          `gorm:"column:blog_id"`
	UserID       uint          `gorm:"column:user_id"`
	ReactionType reaction_type `gorm:"column:reaction_type"`
	PostedOn     time.Time     `gorm:"column:posted_on"`
}

func (*Reaction) TableName() string {
	return "reaction"
}
