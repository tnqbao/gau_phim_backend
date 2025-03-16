package models

type Category struct {
	ID     int     `gorm:"primary_key" json:"id"`
	Name   string  `json:"name"`
	Slug   string  `json:"slug" gorm:"index" gorm:"unique"`
	Movies []Movie `gorm:"many2many:movie_categories;"`
}
