package models

type History struct {
	ID           uint   `json:"id" gorm:"primary_key" `
	UserId       uint   `json:"user_id"`
	MovieName    string `json:"movie_name"`
	MovieSlug    string `json:"movie_slug" gorm:"index;unique"`
	MoviePoster  string `json:"movie_poster"`
	MovieEpisode string `json:"movie_episode"`
	CreateAt     string `json:"create_at"`
}
