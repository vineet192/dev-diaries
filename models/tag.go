package models

import (
	"errors"

	"gorm.io/gorm"
)

type Tag struct {
	ID    uint   `gorm:"primaryKey"`
	Title string `gorm:"column:title"`

	Blogs []Blog `gorm:"many2many:has_tags"`
}

func (t *Tag) validate() (err error) {
	if len(t.Title) == 0 {
		return errors.New("title cannot be empty")
	}
	return
}

func (t *Tag) BeforeSave(tx *gorm.DB) (err error) {
	validateErr := t.validate()
	return validateErr
}

func (*Tag) TableName() string {
	return "tag"
}
