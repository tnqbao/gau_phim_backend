package public

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/config"
	"gorm.io/gorm"
	"strings"
	"time"
)

func GetHomePageData(c *gin.Context) {
	client := config.GetRedisClient()
	db := c.MustGet("db").(*gorm.DB)

	values, _ := client.MGet(ctx,
		"hero_list", "hero_description", "hero_name",
		"featured_homepage_list", "featured_homepage_name", "featured_homepage_year",
	).Result()

	heroSlugs, _ := values[0].(string)
	heroDescriptions, _ := values[1].(string)
	heroNames, _ := values[2].(string)
	featuredSlugs, _ := values[3].(string)
	featuredNames, _ := values[4].(string)
	featuredYears, _ := values[5].(string)

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

	releaseCache, err := client.Get(ctx, "release_list").Result()
	var releaseList []map[string]string
	if err == nil {
		json.Unmarshal([]byte(releaseCache), &releaseList)
	} else {
		var movies []struct {
			Slug string `json:"slug"`
			Name string `json:"name"`
			Year string `json:"year"`
		}
		db.Table("movies").Select("slug, name, year").Limit(24).Scan(&movies)

		for _, movie := range movies {
			releaseList = append(releaseList, map[string]string{
				"slug": movie.Slug,
				"name": movie.Name,
				"year": movie.Year,
			})
		}

		releaseJSON, _ := json.Marshal(releaseList)
		client.Set(ctx, "release_list", releaseJSON, 30*time.Second)
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
