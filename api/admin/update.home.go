package admin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/config"
	"github.com/tnqbao/gau_phim_backend/utils"
	"strings"
)

var ctx = context.Background()

func UpdateHeroHomePage(c *gin.Context) {
	client := config.GetRedisClient()

	var req utils.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": 400, "error": "Invalid request format: " + err.Error()})
		return
	}

	if req.Movies == nil || len(*req.Movies) == 0 {
		c.JSON(400, gin.H{"status": 400, "error": "Movies are required"})
		return
	}

	var heroList, heroDescriptionList []string

	for _, movie := range *req.Movies {
		if movie.Slug != nil {
			heroList = append(heroList, *movie.Slug)
		}
		if movie.Description != nil {
			heroDescriptionList = append(heroDescriptionList, *movie.Description)
		}
	}

	currentHero := strings.Join(heroList, "@")
	currentHeroDescription := strings.Join(heroDescriptionList, "@")

	err := client.MSet(ctx, "hero", currentHero, "hero_description", currentHeroDescription).Err()
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "error": "Failed to update hero list: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Updated hero list successfully"})
}

func UpdateReleaseHomePage(c *gin.Context) {
	client := config.GetRedisClient()

	var req utils.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": 400, "error": "Invalid request format: " + err.Error()})
		return
	}

	if req.Movies == nil || len(*req.Movies) == 0 {
		c.JSON(400, gin.H{"status": 400, "error": "Movies are required"})
		return
	}

	var releaseList []string

	for _, movie := range *req.Movies {
		if movie.Slug != nil {
			releaseList = append(releaseList, *movie.Slug)
		}
	}

	currentRelease := strings.Join(releaseList, "@")

	err := client.Set(ctx, "release_homepage", currentRelease, 0).Err()
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "error": "Failed to update release list: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Updated release list successfully"})
}

func UpdateFeaturedHomePage(c *gin.Context) {
	client := config.GetRedisClient()

	var req utils.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": 400, "error": "Invalid request format: " + err.Error()})
		return
	}

	if req.Movies == nil || len(*req.Movies) == 0 {
		c.JSON(400, gin.H{"status": 400, "error": "Movies are required"})
		return
	}

	var featuredList []string

	for _, movie := range *req.Movies {
		if movie.Slug != nil {
			featuredList = append(featuredList, *movie.Slug)
		}
	}

	currentFeatured := strings.Join(featuredList, "@")

	err := client.Set(ctx, "featured_homepage", currentFeatured, 0).Err()
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "error": "Failed to update featured list: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Updated featured list successfully"})
}
