package authed

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/models"
	"github.com/tnqbao/gau_phim_backend/utils"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func GetHistoryView(c *gin.Context) {
	userId, exists := c.Get("user_id")
	db := c.MustGet("db").(*gorm.DB)
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	if !exists {
		c.JSON(401, gin.H{"error": "User ID not found in context"})
		return
	}

	historyView := []models.History{}
	if err := db.
		Model(&models.History{}).
		Select("slug, title, poster_url, movie_episode, created_at").
		Where("user_id = ?", userId).
		Order("year DESC, modified DESC").
		Limit(24).
		Offset((page - 1) * 24).
		Find(&historyView).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch history view"})
		return
	}

	var totalIteam int64
	if err := db.Model(&models.History{}).
		Where("user_id = ?", userId).
		Count(&totalIteam).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to get total items"})
	}

	totalPage := totalIteam / 24

	c.JSON(500, gin.H{"error": "Fail to get history view"})

	if len(historyView) == 0 {
		c.JSON(404, gin.H{"message": "No history view found"})
		return
	}

	c.JSON(200, gin.H{"history": historyView, "current_page": page, "total_page": totalPage})
}

func UpdateHistoryView(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	var req utils.HistoryRequest
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "User ID not found in context"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if req.MovieName == "" || req.MovieSlug == "" || req.MovieEpisode == "" {
		c.JSON(400, gin.H{"error": "Movie name, slug and episode are required"})
		return
	}

	historyView := models.History{
		UserId:       userId.(uint),
		MovieName:    req.MovieName,
		MovieSlug:    req.MovieSlug,
		MoviePoster:  *req.MoviePoster,
		MovieEpisode: req.MovieEpisode,
		CreateAt:     time.Now().Format("2006-01-02 15:04:05"),
	}
	var existing models.History
	if err := db.Where("user_id = ? AND slug = ?", userId, req.MovieSlug).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&historyView).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to create history view"})
				return
			}
		} else {
			c.JSON(500, gin.H{"error": "DB error when checking history"})
			return
		}
	} else {
		existing.MovieEpisode = req.MovieEpisode
		existing.CreateAt = time.Now().Format("2006-01-02 15:04:05")
		if req.MoviePoster != nil {
			existing.MoviePoster = *req.MoviePoster
		}
		if err := db.Save(&existing).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update history view"})
			return
		}
	}

	c.JSON(200, gin.H{"message": "Update history view successfully"})
}

func DeleteHistoryViewForSlug(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "User ID not found in context"})
		return
	}

	var req utils.HistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if req.MovieSlug == "" {
		c.JSON(400, gin.H{"error": "Movie slug is required"})
		return
	}

	historyView := models.History{
		UserId:    userId.(uint),
		MovieSlug: req.MovieSlug,
	}
	if err := db.Delete(&historyView).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create history view"})
		return
	}

	c.JSON(200, gin.H{"message": "Delete history view successfully"})

}

func DeleteAllHistoryView(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "User ID not found in context"})
		return
	}

	historyView := models.History{
		UserId: userId.(uint),
	}
	if err := db.Delete(&historyView).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create history view"})
		return
	}

	c.JSON(200, gin.H{"message": "Delete all history view successfully"})
}
