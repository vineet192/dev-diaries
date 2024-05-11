package models

type Tag struct {
	id    uint
	title string

	Blogs []Blog `gorm:"many2many:has_tags;foreignKey:tag_id"`
}
