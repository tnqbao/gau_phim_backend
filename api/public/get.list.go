package public

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/config"
	"github.com/tnqbao/gau_phim_backend/models"
	"gorm.io/gorm"
	"strconv"
)

func GetListMovieByType(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	client := config.GetRedisClient()
	slug := c.Param("slug")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	cacheKey := fmt.Sprintf("movies:%s:page:%d", slug, page)

	if cachedData, err := client.Get(ctx, cacheKey).Result(); err == nil {
		client.Expire(ctx, cacheKey, 10*time.Second)

		var cachedResponse map[string]interface{}
		json.Unmarshal([]byte(cachedData), &cachedResponse)
		c.JSON(200, cachedResponse)
		return
	}

	var movies []struct {
		Title     string `json:"title"`
		Year      int    `json:"year"`
		PosterURL string `json:"poster_url"`
		ThumbURL  string `json:"thumb_url"`
		Slug      string `json:"slug"`
	}
	if slug == "new-release" {
		db.Model(&models.Movie{}).
			Select("title, year, poster_url, thumb_url, slug").
			Order("modified desc").
			Limit(24).
			Offset((page - 1) * 24).
			Find(&movies)
	} else {
		db.Model(&models.Movie{}).
			Select("title, year, poster_url, thumb_url, slug").
			Where("type = ?", slug).
			Order("modified desc").
			Limit(24).
			Offset((page - 1) * 24).
			Find(&movies)
	}

	var totalItems int64
	db.Model(&models.Movie{}).Where("type = ?", slug).Count(&totalItems)

	response := gin.H{
		"status":  200,
		"message": "Success",
		"data": gin.H{
			"movies": movies,
			"pagination": gin.H{
				"currentPage":      page,
				"totalItems":       totalItems,
				"totalItemPerPage": 24,
			},
		},
	}

	jsonData, _ := json.Marshal(response)
	client.Set(ctx, cacheKey, jsonData, 10*time.Second)

	c.JSON(200, response)
}
