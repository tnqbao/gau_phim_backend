package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/api/authed/movie"
	"github.com/tnqbao/gau_phim_backend/middlewares"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	apiRoutes := r.Group("/api/gauflix")
	{
		//authedRoutes := forumRoutes.Group("/authed")
		//{
		//	//authedRoutes.Use(middlewares.AuthMiddleware())
		//	//authedRoutes.PUT("/blog", authed.CreateBlog)
		//}
		//publicRoutes := apiRoutes.Group("/")
		//{
		//	//publicRoutes.GET("/blog/:id", public.GetBlogById)
		//}
		adminRoutes := apiRoutes.Group("/admin")
		{
			adminRoutes.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
			adminRoutes.PUT("/crawl", movie.CrawlMovieFromUrl)
		}

	}
	return r
}
