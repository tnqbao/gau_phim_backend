package public

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/config"
	"github.com/tnqbao/gau_phim_backend/models"
	"gorm.io/gorm"
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
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
			c.JSON(200, cachedResponse)
			return
		}
	}

	var movies []struct {
		Title     string `json:"title"`
		Year      int    `json:"year"`
		PosterURL string `json:"poster_url"`
		ThumbURL  string `json:"thumb_url"`
		Slug      string `json:"slug"`
	}

	query := db.Model(&models.Movie{}).Select("title, year, poster_url, thumb_url, slug")

	if slug == "new-release" {
		query = query.Order("year DESC, modified DESC")
	} else {
		query = query.Where("type = ?", slug).Order("year DESC, modified DESC")
	}

	if err := query.Limit(24).Offset((page - 1) * 24).Find(&movies).Error; err != nil {
		c.JSON(500, gin.H{"status": 500, "message": "Database error"})
		return
	}

	var totalItems int64
	countQuery := db.Model(&models.Movie{})
	if slug != "new-release" {
		countQuery = countQuery.Where("type = ?", slug)
	}
	countQuery.Count(&totalItems)

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
