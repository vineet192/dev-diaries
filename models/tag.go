package models

type Tag struct {
	ID    uint   `gorm:"primaryKey;column:id"`
	Title string `gorm:"column:title"`

	Blogs []Blog `gorm:"many2many:has_tags;foreignKey:tag_id"`
}
