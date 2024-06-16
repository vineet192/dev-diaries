package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Blog struct {
	ID       uint      `gorm:"primaryKey;auto_increment"`
	AuthorID uint      `gorm:"column:author_id" json:"author_id"`
	Title    string    `gorm:"column:title" json:"title"`
	Content  string    `gorm:"column:content" json:"content"`
	Markdown bool      `gorm:"column:markdown" json:"markdown"`
	PostedOn time.Time `gorm:"column:posted_on" json:"posted_on"`

	Tags      []Tag          `gorm:"many2many:has_tags;constraint:onDelete:CASCADE" json:"tags"`
	Reactions []BlogReaction `gorm:"foreignKey:blog_id;constraint:onDelete:CASCADE" json:"reactions"`
	Comments  []Comment      `gorm:"foreignKey:blog_id;constraint:onDelete:CASCADE" json:"comments"`
}

func (*Blog) TableName() string {
	return "blog"
}

func (b *Blog) validate() (err error) {

	if b.Title == "" {
		return errors.New("blog title cannot be empty")
	}

	return nil
}

func (b *Blog) BeforeSave(tx *gorm.DB) (err error) {
	validateErr := b.validate()
	b.PostedOn = time.Now()
	return validateErr
}
