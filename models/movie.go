package models

import "time"

type Movie struct {
	ID          int        `gorm:"primary_key" json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Year        int        `json:"year"`
	PosterURL   string     `json:"poster_url"`
	ThumbUrl    string     `json:"thumb_url"`
	Slug        string     `json:"slug" gorm:"index;unique"`
	Nations     []Nation   `gorm:"many2many:movie_nations;"`
	Categories  []Category `gorm:"many2many:movie_categories;"`
	Type        string     `json:"type"`
	Modified    time.Time  `json:"modified"`
}
