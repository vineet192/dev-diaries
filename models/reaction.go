package models

import (
	"time"

	"gorm.io/gorm"
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

func (r *Reaction) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Model(r).Update("posted_on", time.Now())
	return
}
