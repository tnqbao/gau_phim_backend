package utils

import "github.com/tnqbao/gau_phim_backend/models"

type Request struct {
	Slug        *string            `json:"slug"`
	Page        *int               `json:"page"`
	Endpoint    *string            `json:"endpoint"`
	Amount      *int               `json:"amount"`
	Movie       *models.Movie      `json:"movie"`
	PosterUrl   *string            `json:"poster_url"`
	ThumbUrl    *string            `json:"thumb_url"`
	Title       *string            `json:"title"`
	Year        *int               `json:"year"`
	Description *string            `json:"description"`
	OriginTitle *string            `json:"origin_title"`
	Nations     *[]models.Nation   `json:"nations"`
	Categories  *[]models.Category `json:"categories"`
}

type ApiResponse struct {
	Data struct {
		Items []struct {
			Name       string `json:"name"`
			OriginName string `json:"origin_name"`
			Slug       string `json:"slug"`
			Type       string `json:"type"`
			ThumbURL   string `json:"thumb_url"`
			PosterURL  string `json:"poster_url"`
			Time       string `json:"time"`
			Year       int    `json:"year"`
			Quality    string `json:"quality"`
			Lang       string `json:"lang"`
			Categories []struct {
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"category"`
			Countries []struct {
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"country"`
		} `json:"items"`
	} `json:"data"`
}
