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

func V1UserInsertPost(c *gin.Context, blogData db.BlogDataV1, userToken string) {
	valid, email := V1VerifyToken(userToken)
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "invalid token",
		})
		return
	}
	// check user permission
	if email != "admin" {
		_ = db.V1InsertPost(blogData)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
	}
}
