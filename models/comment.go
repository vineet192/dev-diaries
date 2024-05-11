package models

import "time"

type Comment struct {
	id        uint
	user_id   uint
	blog_id   uint
	content   string
	posted_on time.Time
	last_edit time.Time

	Reactions []Reaction `gorm:"foreignKey:comment_id"`
}
