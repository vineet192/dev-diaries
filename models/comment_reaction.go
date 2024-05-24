package models

import (
	"time"

	"gorm.io/gorm"
)

type CommentReaction struct {
	CommentID    uint          `gorm:"primaryKey;column:comment_id;default:null" json:"comment_id"`
	UserID       uint          `gorm:"primaryKey;column:user_id" json:"user_id"`
	ReactionType reaction_type `gorm:"column:reaction_type" json:"reaction_type"`
	PostedOn     time.Time     `gorm:"column:posted_on" json:"posted_on"`
}

func (*CommentReaction) TableName() string {
	return "comment_reaction"
}

func (r *CommentReaction) BeforeCreate(tx *gorm.DB) (err error) {
	r.PostedOn = time.Now()
	return
}
