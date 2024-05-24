package models

import (
	"time"

	"gorm.io/gorm"
)

type reaction_type string

type BlogReaction struct {
	BlogID       uint          `gorm:"primaryKey;column:blog_id;default:null" json:"blog_id"`
	UserID       uint          `gorm:"primaryKey;column:user_id" json:"user_id"`
	ReactionType reaction_type `gorm:"column:reaction_type" json:"reaction_type"`
	PostedOn     time.Time     `gorm:"column:posted_on" json:"posted_on"`
}

func (*BlogReaction) TableName() string {
	return "blog_reaction"
}

func (r *BlogReaction) BeforeCreate(tx *gorm.DB) (err error) {
	r.PostedOn = time.Now()
	return
}
