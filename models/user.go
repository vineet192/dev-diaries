package models

import (
	"errors"
	"regexp"

	"gorm.io/gorm"
)

type User struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string     `gorm:"column:first_name" json:"first_name"`
	LastName  string     `gorm:"column:last_name" json:"last_name"`
	Email     string     `gorm:"column:email" json:"email"`
	Bio       string     `gorm:"column:bio" json:"bio"`
	Blogs     []Blog     `gorm:"foreignKey:author_id"`
	Reactions []Reaction `gorm:"foreignKey:user_id"`
	Comments  []Comment  `gorm:"foreignKey:user_id"`
}

func (*User) TableName() string {
	return "user"
}

func (u *User) validate() (err error) {

	if u.FirstName == "" {
		return errors.New("first name cannot be empty")
	}

	emailMatched, _ := regexp.MatchString("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", u.Email)

	if !emailMatched {
		return errors.New("invalid email format")
	}

	if len(u.Bio) > 200 {
		return errors.New("bio too long")
	}

	return
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	validationErr := u.validate()
	return validationErr
}
