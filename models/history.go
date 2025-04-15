package models

type History struct {
	ID           uint   `json:"id" gorm:"primary_key" `
	UserId       uint   `json:"-"`
	MovieName    string `json:"title"`
	MovieSlug    string `json:"slug" gorm:"index;unique"`
	MoviePoster  string `json:"poster_url"`
	MovieEpisode string `json:"movie_episode"`
	CreateAt     string `json:"create_at"`
}
