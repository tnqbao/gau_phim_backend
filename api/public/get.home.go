package public

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/config"
	"strings"
)

var ctx = context.Background()

func GetHomePageData(c *gin.Context) {
	client := config.GetRedisClient()

	heroSlugs, _ := client.Get(ctx, "hero").Result()
	heroDescriptions, _ := client.Get(ctx, "hero_description").Result()
	releaseSlugs, _ := client.Get(ctx, "release_homepage").Result()
	featuredSlugs, _ := client.Get(ctx, "featured_homepage").Result()

	heroList := []map[string]string{}
	heroSlugArr := strings.Split(heroSlugs, "@")
	heroDescriptionArr := strings.Split(heroDescriptions, "@")

	for i := range heroSlugArr {
		if i < len(heroDescriptionArr) {
			heroList = append(heroList, map[string]string{
				"slug":        heroSlugArr[i],
				"description": heroDescriptionArr[i],
			})
		}
	}

	releaseList := strings.Split(releaseSlugs, "@")
	featuredList := strings.Split(featuredSlugs, "@")

	c.JSON(200, gin.H{
		"hero":     heroList,
		"release":  releaseList,
		"featured": featuredList,
	})
}
