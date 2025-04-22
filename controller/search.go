package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/meilisearch/meilisearch-go"
	"github.com/tnqbao/gau_phim_backend/config"
	"github.com/tnqbao/gau_phim_backend/models"
	"net/http"
	"strings"
)

func SearchMovieByKeyWord(c *gin.Context) {
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keyword is required"})
		return
	}

	index := config.MeiliClient.Index("movies")

	searchRequest := &meilisearch.SearchRequest{
		Limit:                 10,
		Offset:                0,
		AttributesToHighlight: []string{"title"},
		AttributesToCrop:      []string{"title"},
		HighlightPreTag:       "<strong>",
		HighlightPostTag:      "</strong>",
	}

	searchResponse, err := index.Search(keyword, searchRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search movies", "detail": err.Error()})
		return
	}

	results := make([]gin.H, 0)
	for _, hit := range searchResponse.Hits {
		h := hit.(map[string]interface{})
		results = append(results, gin.H{
			"id":         h["id"],
			"slug":       h["slug"],
			"title":      h["title"],
			"year":       h["year"],
			"_highlight": h["_formatted"],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"hits":  results,
		"total": searchResponse.EstimatedTotalHits,
	})
}

func IndexAllMovies(c *gin.Context) {
	var movies []models.Movie
	if err := config.DB.Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies from database"})
		return
	}

	var movieIndexes []models.MovieIndex
	for _, m := range movies {
		movieIndexes = append(movieIndexes, models.MovieIndex{
			ID:    m.ID,
			Slug:  m.Slug,
			Title: m.Title,
			Year:  m.Year,
		})
	}

	index := config.MeiliClient.Index("movies")

	_, err := index.UpdateSearchableAttributes(&[]string{"title"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set searchable attributes"})
		return
	}

	_, err = index.AddDocuments(movieIndexes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to index movies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ Movies indexed successfully"})
}
