package models

type Tag struct {
	ID    uint   `gorm:"primaryKey"`
	Title string `gorm:"column:title"`

	Blogs []Blog `gorm:"many2many:has_tags"`
}

func (*Tag) TableName() string {
	return "tag"
}
