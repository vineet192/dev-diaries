package models

import "time"

type Comment struct {
	ID       uint      `gorm:"primaryKey;column:id"`
	UserID   uint      `gorm:"column:user_id"`
	BlogID   uint      `gorm:"column:blog_id"`
	Content  string    `gorm:"column:content"`
	PostedOn time.Time `gorm:"column:posted_on"`
	LastEdit time.Time `gorm:"column:last_edit"`

	Reactions []Reaction `gorm:"foreignKey:comment_id"`
}
