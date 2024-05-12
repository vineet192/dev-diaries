package models

type User struct {
	ID        uint       `gorm:"primaryKey;column:id"`
	FirstName string     `gorm:"column:first_name"`
	LastName  string     `gorm:"column:last_name"`
	Email     string     `gorm:"column:email"`
	Bio       string     `gorm:"column:bio"`
	Blogs     []Blog     `gorm:"foreignKey:author_id`
	Reactions []Reaction `gorm:"foreignKey:user_id"`
	Comments  []Comment  `gorm:"foreignKey:user_id"`
}
