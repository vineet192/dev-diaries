package models

import (
	"errors"
	"regexp"

	"gorm.io/gorm"
)

type User struct {
	ID               uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName        string            `gorm:"column:first_name" json:"first_name"`
	LastName         string            `gorm:"column:last_name" json:"last_name"`
	Email            string            `gorm:"column:email;unique" json:"email"`
	Bio              string            `gorm:"column:bio" json:"bio"`
	Hash             string            `gorm:"column:hash;not null;"`
	Blogs            []Blog            `gorm:"foreignKey:AuthorID;constraint:onDelete:CASCADE"`
	BlogReactions    []BlogReaction    `gorm:"foreignKey:user_id;constraint:onDelete:CASCADE"`
	CommentReactions []CommentReaction `gorm:"foreignKey:user_id;constraint:onDelete:CASCADE"`
	Comments         []Comment         `gorm:"foreignKey:user_id;constraint:onDelete:CASCADE"`
	Followers        []User            `gorm:"many2many:has_followers;joinForeignKey:UserID;joinReferences:FollowerID;constraint:onDelete:CASCADE;"`
	Following        []User            `gorm:"many2many:has_followers;joinForeignKey:FollowerID;joinReferences:UserID;constraint:onDelete:CASCADE;"`
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
