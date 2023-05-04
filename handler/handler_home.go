package handler
import (
    "github.com/gin-gonic/gin"
    "net/http"
    db "go_blog/database"
)

func HandlerHome(c *gin.Context) {
    // get 20 recent posts
    posts, _ := db.GetRecentPosts(20)
    c.JSON(http.StatusOK, posts)
}


