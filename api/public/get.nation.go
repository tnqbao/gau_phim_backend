package public

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/config"
	"github.com/tnqbao/gau_phim_backend/models"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func GetListMovieByNation(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	client := config.GetRedisClient()
	slug := c.Param("slug")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	cacheKey := "movies:nation:" + slug + ":page:" + strconv.Itoa(page)

	if cachedData, err := client.Get(ctx, cacheKey).Result(); err == nil {
		client.Expire(ctx, cacheKey, 30*time.Second)
		c.JSON(200, gin.H{
			"status":  200,
			"message": "Success (from cache)",
			"data":    json.RawMessage(cachedData),
		})
		return
	}

	var movies []struct {
		Title     string `json:"title"`
		Year      int    `json:"year"`
		PosterURL string `json:"poster_url"`
		ThumbURL  string `json:"thumb_url"`
		Slug      string `json:"slug"`
	}

	db.Model(&models.Movie{}).
		Select("movies.title, movies.year, movies.poster_url, movies.thumb_url, movies.slug").
		Joins("JOIN movie_nations mn ON movies.id = mn.movie_id").
		Joins("JOIN nations n ON mn.nation_id = n.id").
		Where("n.slug = ?", slug).
		Order("movies.year DESC, movies.modified DESC").
		Limit(24).
		Offset((page - 1) * 24).
		Find(&movies)

	var totalItems int64
	db.Model(&models.Movie{}).
		Joins("JOIN movie_nations mn ON movies.id = mn.movie_id").
		Joins("JOIN nations n ON mn.nation_id = n.id").
		Where("n.slug = ?", slug).
		Count(&totalItems)

	responseData := gin.H{
		"movies": movies,
		"pagination": gin.H{
			"currentPage":      page,
			"totalItems":       totalItems,
			"totalItemPerPage": 24,
		},
	}

	jsonData, _ := json.Marshal(responseData)

	client.Set(ctx, cacheKey, jsonData, 30*time.Second)

	c.JSON(200, gin.H{
		"status":  200,
		"message": "Success",
		"data":    responseData,
	})
}
