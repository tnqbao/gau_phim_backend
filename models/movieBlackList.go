package models

type MovieBlackList struct {
	ID   int    `gorm:"primary_key" json:"id"`
	Slug string `json:"slug" gorm:"index;unique"`
}
