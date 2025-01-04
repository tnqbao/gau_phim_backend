package routes

// func SetupRouter(db *gorm.DB) *gin.Engine {
// 	r := gin.Default()
// 	r.Use(middlewares.CORSMiddleware())
// 	r.Use(func(c *gin.Context) {
// 		c.Set("db", db)
// 		c.Next()
// 	})
//  *** write your routes below ***
//
//  ** Example :
// 	apiRoutes := r.Group("/api")
// 	{
// 		subRoutes := apiRoutes.Group("/your-path")
// 		{
// 			authedRoutes := forumRoutes.Group("/authed")
// 			{
// 				authedRoutes.Use(middlewares.AuthMiddleware())
// 				authedRoutes.PUT("/blog", authed.CreateBlog)
// 			}
// 			publicRoutes := forumRoutes.Group("/public")
// 			{
// 				publicRoutes.GET("/blog/:id", public.GetBlogById)
// 			}
// 		}
// 	}
//  **
// 	return r
// }
