package admin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/config"
	"github.com/tnqbao/gau_phim_backend/utils"
	"strings"
)

var ctx = context.Background()

func updateHeroHomePage(c *gin.Context) {
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

	var slugList, nameList, descriptionList []string

	for _, movie := range *req.Movies {
		if movie.Slug != nil {
			slugList = append(slugList, *movie.Slug)
		}
		if movie.Name != nil {
			nameList = append(nameList, *movie.Name)
		}
		if movie.Description != nil {
			descriptionList = append(descriptionList, *movie.Description)
		}
	}

	pipe := client.Pipeline()
	pipe.Set(ctx, "hero_list", strings.Join(slugList, "@"), 0)
	pipe.Set(ctx, "hero_name", strings.Join(nameList, "@"), 0)
	if len(descriptionList) > 0 {
		pipe.Set(ctx, "hero_description", strings.Join(descriptionList, "@"), 0)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "error": "Failed to update hero list: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Updated hero list successfully"})
}

func updateGeneralHomePage(c *gin.Context, keyPrefix string) {
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

	var slugList, nameList, yearList []string

	for _, movie := range *req.Movies {
		if movie.Slug != nil {
			slugList = append(slugList, *movie.Slug)
		}
		if movie.Name != nil {
			nameList = append(nameList, *movie.Name)
		}
		if movie.Year != nil {
			yearList = append(yearList, *movie.Year)
		}
	}

	pipe := client.Pipeline()
	pipe.Set(ctx, keyPrefix+"_list", strings.Join(slugList, "@"), 0)
	pipe.Set(ctx, keyPrefix+"_name", strings.Join(nameList, "@"), 0)
	if len(yearList) > 0 {
		pipe.Set(ctx, keyPrefix+"_year", strings.Join(yearList, "@"), 0)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "error": "Failed to update " + keyPrefix + " list: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Updated " + keyPrefix + " list successfully"})
}

func UpdateHeroHomePage(c *gin.Context) {
	updateHeroHomePage(c)
}

func UpdateReleaseHomePage(c *gin.Context) {
	updateGeneralHomePage(c, "release_homepage")
}

func UpdateFeaturedHomePage(c *gin.Context) {
	updateGeneralHomePage(c, "featured_homepage")
}
