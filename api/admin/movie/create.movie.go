package movie

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/models"
	"github.com/tnqbao/gau_phim_backend/utils"
	"gorm.io/gorm"
	"log"
	"math"
	"net/http"
	"net/url"
	"time"
)

func CreateMovie(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Create movie",
	})
}

func CrawlMovieFromUrl(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	params := url.Values{}
	var req utils.Request

	count := 0
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Invalid request format: " + err.Error(),
		})
		return
	}

	if req.Endpoint == nil || req.Amount == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Endpoint and Amount are required",
		})
		return
	}

	amountPage := int(math.Ceil(float64(*req.Amount) / 24))

	for i := 1; i <= amountPage; i++ {
		params.Set("page", fmt.Sprintf("%d", i))
		url := fmt.Sprintf("%s?%s", *req.Endpoint, params.Encode())
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Lỗi khi gọi API trang %d: %v", i, err)
			continue
		}

		defer resp.Body.Close()

		var apiResp utils.ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			log.Printf("Lỗi khi decode JSON trang %d: %v", i, err)
			continue
		}

		for _, item := range apiResp.Data.Items {
			var existingMovie models.Movie
			if err := db.Where("slug = ?", item.Slug).First(&existingMovie).Error; err == nil {
				log.Printf("Bỏ qua: Phim %s đã tồn tại", item.Name)
				continue
			}

			var categories []models.Category
			for _, cat := range item.Categories {
				var category models.Category
				if err := db.Where("slug = ?", cat.Slug).FirstOrCreate(&category, models.Category{
					Name: cat.Name, Slug: cat.Slug,
				}).Error; err != nil {
					log.Printf("Lỗi khi thêm thể loại %s: %v", cat.Name, err)
				}
				categories = append(categories, category)
			}

			var countries []models.Nation
			for _, country := range item.Countries {
				var nation models.Nation
				if err := db.Where("slug = ?", country.Slug).FirstOrCreate(&nation, models.Nation{
					Name: country.Name, Slug: country.Slug,
				}).Error; err != nil {
					log.Printf("Lỗi khi thêm quốc gia %s: %v", country.Name, err)
				}
				countries = append(countries, nation)
			}

			// Parse modified time
			modifiedTime, err := time.Parse(time.RFC3339, item.Modified.Time)
			if err != nil {
				log.Printf("Lỗi khi parse thời gian %s: %v", item.Modified.Time, err)
				modifiedTime = time.Now()
			}

			movie := models.Movie{
				Title:      item.Name,
				Slug:       item.Slug,
				Year:       item.Year,
				PosterURL:  item.PosterURL,
				ThumbUrl:   item.ThumbURL,
				Categories: categories,
				Nations:    countries,
				Type:       item.Type,
				Modified:   modifiedTime,
			}

			if err := db.Create(&movie).Error; err != nil {
				log.Printf("Lỗi khi lưu phim %s: %v", movie.Title, err)
			} else {
				count++
				log.Printf("Đã lưu phim: %s", movie.Title)
			}
		}

		fmt.Printf("Đã nhận phản hồi cho trang %d: %d\n", i, resp.StatusCode)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Crawl movie from URL hoàn tất",
		"Đã thêm": count,
	})
}
