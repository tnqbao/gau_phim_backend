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

	values, _ := client.MGet(ctx,
		"hero_list", "hero_description", "hero_name",
		"release_homepage_list", "release_homepage_name", "release_homepage_year",
		"featured_homepage_homepage", "featured_homepage_name", "featured_homepage_year",
	).Result()

	heroSlugs, _ := values[0].(string)
	heroDescriptions, _ := values[1].(string)
	heroNames, _ := values[2].(string)
	releaseSlugs, _ := values[3].(string)
	releaseNames, _ := values[4].(string)
	releaseYears, _ := values[5].(string)
	featuredSlugs, _ := values[6].(string)
	featuredNames, _ := values[7].(string)
	featuredYears, _ := values[8].(string)

	heroList := []map[string]string{}
	heroSlugArr := strings.Split(heroSlugs, "@")
	heroDescriptionArr := strings.Split(heroDescriptions, "@")
	heroNameArr := strings.Split(heroNames, "@")
	for i := range heroSlugArr {
		if i < len(heroDescriptionArr) {
			heroList = append(heroList, map[string]string{
				"slug":        heroSlugArr[i],
				"name":        heroNameArr[i],
				"description": heroDescriptionArr[i],
			})
		}
	}

	releaseList := []map[string]string{}
	releaseSlugArr := strings.Split(releaseSlugs, "@")
	releaseNameArr := strings.Split(releaseNames, "@")
	releaseYearArr := strings.Split(releaseYears, "@")
	for i := range releaseSlugArr {
		if i < len(releaseNameArr) && i < len(releaseYearArr) {
			releaseList = append(releaseList, map[string]string{
				"slug": releaseSlugArr[i],
				"name": releaseNameArr[i],
				"year": releaseYearArr[i],
			})
		}
	}

	featuredList := []map[string]string{}
	featuredSlugArr := strings.Split(featuredSlugs, "@")
	featuredNameArr := strings.Split(featuredNames, "@")
	featuredYearArr := strings.Split(featuredYears, "@")
	for i := range featuredSlugArr {
		if i < len(featuredNameArr) && i < len(featuredYearArr) {
			featuredList = append(featuredList, map[string]string{
				"slug": featuredSlugArr[i],
				"name": featuredNameArr[i],
				"year": featuredYearArr[i],
			})
		}
	}

	c.JSON(200, gin.H{
		"hero":     heroList,
		"release":  releaseList,
		"featured": featuredList,
	})
}
