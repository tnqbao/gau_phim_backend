package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/config"
	"github.com/tnqbao/gau_phim_backend/models"
	"github.com/tnqbao/gau_phim_backend/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	index := config.MeiliClient.Index("movies")
	_, _ = index.UpdateSearchableAttributes(&[]string{"title"})

	err := db.Transaction(func(tx *gorm.DB) error {
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
				var black models.MovieBlackList
				if err := db.Where("slug = ?", item.Slug).First(&black).Error; err == nil {
					log.Printf("🚫 Bỏ qua (blacklist): %s", item.Slug)
					continue
				}

				var existingMovie models.Movie
				if err := db.Where("slug = ?", item.Slug).First(&existingMovie).Error; err == nil {
					log.Printf("Bỏ qua: Phim %s đã tồn tại", item.Name)
					continue
				}

				var categories []models.Category
				for _, cat := range item.Categories {
					var category models.Category
					if err := tx.Where("slug = ?", cat.Slug).FirstOrCreate(&category, models.Category{
						Name: cat.Name, Slug: cat.Slug,
					}).Error; err != nil {
						log.Printf("Lỗi khi thêm thể loại %s: %v", cat.Name, err)
					}
					categories = append(categories, category)
				}

				var countries []models.Nation
				for _, country := range item.Countries {
					var nation models.Nation
					if err := tx.Where("slug = ?", country.Slug).FirstOrCreate(&nation, models.Nation{
						Name: country.Name, Slug: country.Slug,
					}).Error; err != nil {
						log.Printf("Lỗi khi thêm quốc gia %s: %v", country.Name, err)
					}
					countries = append(countries, nation)
				}

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

				if err := tx.Create(&movie).Error; err != nil {
					log.Printf("Lỗi khi lưu phim %s: %v", movie.Title, err)
				} else {
					count++
					log.Printf("✅ Đã lưu phim: %s", movie.Title)
					movieToIndex := models.MovieIndex{
						ID:    movie.ID,
						Slug:  movie.Slug,
						Title: movie.Title,
						Year:  movie.Year,
					}
					_, err := index.AddDocuments([]models.MovieIndex{movieToIndex})
					if err != nil {
						log.Printf("❌ Lỗi khi index phim %s: %v", movie.Title, err)
					} else {
						log.Printf("🔍 Đã index phim: %s", movie.Title)
					}
				}
			}
			fmt.Printf("Đã nhận phản hồi cho trang %d: %d\n", i, resp.StatusCode)
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Lỗi khi crawl phim: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Crawl movie from URL hoàn tất",
		"Đã thêm": count,
	})
}

func GetAllMovie(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var movies []models.Movie
	if err := db.Preload("Categories").Preload("Nations").Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Lỗi khi lấy danh sách phim: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   movies,
	})
}

func DeleteMovieBySlug(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	slug := c.Param("slug")

	err := db.Transaction(func(tx *gorm.DB) error {
		var movie models.Movie
		if err := tx.Where("slug = ?", slug).First(&movie).Error; err != nil {
			return err
		}

		if err := tx.Delete(&movie).Error; err != nil {
			return err
		}

		index := config.MeiliClient.Index("movies")
		if deleteRes, err := index.DeleteDocument(slug); err != nil {
			log.Printf("❌ Lỗi khi xoá index Meilisearch cho slug %s: %v", slug, err)
		} else {
			log.Printf("✅ Đã xóa index Meilisearch cho slug %s: %v", slug, deleteRes)
		}

		// Thêm vào bảng blacklist
		black := models.MovieBlackList{Slug: slug}
		if err := tx.FirstOrCreate(&black, models.MovieBlackList{Slug: slug}).Error; err != nil {
			return fmt.Errorf("xóa thành công nhưng lỗi khi thêm vào blacklist: %w", err)
		}

		return nil
	})

	if err != nil {
		status := http.StatusInternalServerError
		msg := "Lỗi khi xóa phim: " + err.Error()
		if err == gorm.ErrRecordNotFound {
			status = http.StatusNotFound
			msg = "Phim không tồn tại"
		}
		c.JSON(status, gin.H{
			"status": status,
			"error":  msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Đã xóa phim và thêm vào blacklist: %s", slug),
	})
}

func DeleteMovieByListSlug(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var slugs []string
	if err := c.ShouldBindJSON(&slugs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Lỗi khi phân tích cú pháp danh sách slug: " + err.Error(),
		})
		return
	}

	if len(slugs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Danh sách slug không được để trống",
		})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("slug IN ?", slugs).Delete(&models.Movie{}).Error; err != nil {
			return err
		}

		index := config.MeiliClient.Index("movies")
		if deleteResp, err := index.DeleteDocuments(slugs); err != nil {
			log.Printf("❌ Lỗi khi xoá index Meilisearch: %v", err)
		} else {
			log.Printf("✅ Đã xóa index Meilisearch: %v", deleteResp)
		}

		// Thêm tất cả slug vào bảng blacklist
		var blacklists []models.MovieBlackList
		for _, slug := range slugs {
			blacklists = append(blacklists, models.MovieBlackList{Slug: slug})
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&blacklists).Error; err != nil {
			return fmt.Errorf("đã xóa phim nhưng lỗi khi thêm vào blacklist: %w", err)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Lỗi khi xóa phim: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Đã xóa %d phim và thêm vào blacklist", len(slugs)),
	})
}
