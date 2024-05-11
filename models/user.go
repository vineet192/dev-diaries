package models

type User struct {
	id         uint `gorm:"primaryKey"`
	first_name string
	last_name  string
	email      string
	bio        string

	Blogs     []Blog     `gorm:"foreignKey:author_id`
	Reactions []Reaction `gorm:"foreignKey:user_id"`
	Comments  []Comment  `gorm:"foreignKey:user_id"`
}
