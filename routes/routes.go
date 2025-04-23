package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_phim_backend/controller"
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
		publicRouter := apiRoutes.Group("/")
		{
			publicRouter.GET("/all-movies", controller.GetAllMovie)
			publicRouter.GET("/home-page", controller.GetHomePageData)
			publicRouter.GET("/category/:slug", controller.GetListMovieByCategory)
			publicRouter.GET("/type/:slug", controller.GetListMovieByType)
			publicRouter.GET("/nation/:slug", controller.GetListMovieByNation)

			publicRouter.POST("/search", controller.SearchMovieByKeyWord)

		}
		adminRoutes := apiRoutes.Group("/")
		{
			adminRoutes.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
			adminRoutes.PUT("/crawl", controller.CrawlMovieFromUrl)
			adminRoutes.POST("/movie", controller.CreateMovie)
			adminRoutes.DELETE("/movie/:slug", controller.DeleteMovieBySlug)
			adminRoutes.POST("/movie", controller.DeleteMovieByListSlug)

			adminRoutes.PUT("/home-page/hero", controller.UpdateHeroHomePage)
			adminRoutes.PUT("/home-page/release", controller.UpdateReleaseHomePage)
			adminRoutes.PUT("/home-page/featured", controller.UpdateFeaturedHomePage)
			//search
			adminRoutes.POST("/index", controller.IndexAllMovies)
		}

		authedRoutes := apiRoutes.Group("/")
		{
			authedRoutes.Use(middlewares.AuthMiddleware())
			authedRoutes.POST("/like", controller.AddMovieLiked)
			authedRoutes.GET("/likes", controller.GetListMovieLiked)
			authedRoutes.DELETE("/like", controller.RemoveMovieLiked)

			authedRoutes.GET("/history", controller.GetHistoryWatched)
			authedRoutes.POST("/history", controller.UpdateHistoryWatched)
			authedRoutes.DELETE("/history/:slug", controller.DeleteHistoryWatchedForSlug)
			authedRoutes.DELETE("/history", controller.DeleteAllHistoryWatched)
		}

	}
	return r
}
