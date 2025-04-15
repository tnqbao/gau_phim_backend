package authed

//
//import "github.com/gin-gonic/gin"
//
//func AddPlayList(c *gin.Context) {
//	// Get the user ID from the context
//	userId, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(401, gin.H{"error": "User ID not found in context"})
//		return
//	}
//
//	// Get the playlist name from the request body
//	var requestBody struct {
//		Title   string `json:"title" binding:"required"`
//		Episode string `json:"episode" binding:"required"`
//	}
//	if err := c.ShouldBindJSON(&requestBody); err != nil {
//		c.JSON(400, gin.H{"error": "Invalid request body"})
//		return
//	}
//
//	// Create a new playlist in the database
//	playList := models.PlayList{
//		UserId: userId.(uint),
//		Name:   requestBody.Name,
//	}
//	if err := models.CreatePlayList(&playList); err != nil {
//		c.JSON(500, gin.H{"error": "Failed to create playlist"})
//		return
//	}
//
//	c.JSON(200, playList)
//}
