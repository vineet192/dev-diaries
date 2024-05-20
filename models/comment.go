package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID       uint      `gorm:"primaryKey"`
	UserID   uint      `gorm:"column:user_id"`
	BlogID   uint      `gorm:"column:blog_id"`
	Content  string    `gorm:"column:content"`
	PostedOn time.Time `gorm:"column:posted_on"`
	LastEdit time.Time `gorm:"column:last_edit"`

	Reactions []Reaction `gorm:"foreignKey:comment_id"`
}

func (*Comment) TableName() string {
	return "comment"
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Model(c).Update("posted_on", time.Now())
	tx.Model(c).Update("last_edit", time.Now())
	return
}

func (c *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Model(c).Update("last_edit", time.Now())
	return
}
