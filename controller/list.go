package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/config"
	"github.com/tnqbao/gau_phim_backend/models"
	"gorm.io/gorm"
	"strconv"
	"time"
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

func GetListMovieByCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	client := config.GetRedisClient()
	slug := c.Param("slug")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	cacheKey := "movies:category:" + slug + ":page:" + strconv.Itoa(page)

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
		Joins("JOIN movie_categories mc ON movies.id = mc.movie_id").
		Joins("JOIN categories c ON mc.category_id = c.id").
		Where("c.slug = ?", slug).
		Order("movies.year DESC, movies.modified DESC").
		Limit(24).
		Offset((page - 1) * 24).
		Find(&movies)

	var totalItems int64
	db.Model(&models.Movie{}).
		Joins("JOIN movie_categories mc ON movies.id = mc.movie_id").
		Joins("JOIN categories c ON mc.category_id = c.id").
		Where("c.slug = ?", slug).
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
