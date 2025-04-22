package models

type MovieIndex struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Year  int    `json:"year"`
}
