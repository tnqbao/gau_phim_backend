package authed

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/models"
	"github.com/tnqbao/gau_phim_backend/utils"
	"gorm.io/gorm"
)

func AddMovieLiked(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.MustGet("user_id")
	if id == nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	userID := id.(int)

	var req utils.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var movie models.Movie
	if err := db.First(&movie, "slug = ?", req.Slug).Error; err != nil {
		c.JSON(404, gin.H{"error": "Movie not found"})
		return
	}

	var like models.MovieLike
	if err := db.Where("user_id = ? AND movie_id = ?", userID, movie.ID).First(&like).Error; err == nil {
		c.JSON(400, gin.H{"error": "You already liked this movie"})
		return
	}

	like = models.MovieLike{
		UserID:  userID,
		MovieID: movie.ID,
	}

	if err := db.Create(&like).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to add like"})
		return
	}

	c.JSON(200, gin.H{
		"status":  200,
		"message": "Liked movie successfully",
	})
}

func RemoveMovieLiked(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.MustGet("user_id")
	if id == nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	userID := id.(int)

	var req utils.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var movie models.Movie
	if err := db.First(&movie, "slug = ?", req.Slug).Error; err != nil {
		c.JSON(404, gin.H{"error": "Movie not found"})
		return
	}

	var like models.MovieLike
	if err := db.Where("user_id = ? AND movie_id = ?", userID, movie.ID).First(&like).Error; err != nil {
		c.JSON(400, gin.H{"error": "You have not liked this movie"})
		return
	}

	if err := db.Delete(&like).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to remove like"})
		return
	}

	c.JSON(200, gin.H{
		"status":  200,
		"message": "Removed like successfully",
	})
}
func GetListMovieLiked(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.MustGet("user_id")
	if id == nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	userID := id.(uint)

	var likes []models.MovieLike
	if err := db.Where("user_id = ?", userID).Find(&likes).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to get liked movies"})
		return
	}

	var movies []models.Movie
	for _, like := range likes {
		var movie models.Movie
		if err := db.First(&movie, like.MovieID).Error; err != nil {
			continue
		}
		movies = append(movies, movie)
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   movies,
	})
}
