package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID       uint      `gorm:"primaryKey"`
	UserID   uint      `gorm:"column:user_id" json:"user_id"`
	BlogID   uint      `gorm:"column:blog_id" json:"blog_id"`
	Content  string    `gorm:"column:content" json:"content"`
	PostedOn time.Time `gorm:"column:posted_on" json:"posted_on"`
	LastEdit time.Time `gorm:"column:last_edit" json:"last_edit"`

	Reactions []Reaction `gorm:"foreignKey:comment_id"`
}

func (*Comment) TableName() string {
	return "comment"
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.PostedOn = time.Now()
	c.LastEdit = time.Now()
	return
}

func (c *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	c.LastEdit = time.Now()
	return
}
