package models

type History struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	UserId       uint   `json:"-"`
	MovieName    string `json:"title" gorm:"column:title"`
	MovieSlug    string `json:"slug" gorm:"index;column:slug"`
	MoviePoster  string `json:"poster_url" gorm:"column:poster_url"`
	MovieEpisode string `json:"movie_episode" gorm:"column:movie_episode"`
	CreateAt     string `json:"create_at" gorm:"column:create_at"`
}
