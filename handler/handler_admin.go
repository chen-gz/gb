package handler

import (
	"github.com/gin-gonic/gin"
	db "go_blog/database"
	rd "go_blog/render"
	"net/http"
)

func V1UserGetPost(c *gin.Context, url string, userEmail string) {
	// check user permission
	post := db.V1GetPostByUrl(url)
	user := db.GetUser(userEmail)
	if user.Group == "admin" || user.Level >= post.PrivateLevel {
		post.Rendered = string(rd.RenderMd([]byte(post.Content)))
		c.JSON(http.StatusOK, post)
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "permission denied",
		})
	}
}

func V1UserUpdatePost(c *gin.Context, blogData db.BlogDataV1, userEmail string) {
	// check user permission
	// todo: this should get by id instead of url
	user := db.GetUser(userEmail)

	if user.Group == "admin" {
		db.V1UpdatePost(blogData)
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "permission denied",
		})
	}
}
