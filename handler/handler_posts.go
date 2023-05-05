package handler

import (
	"github.com/gin-gonic/gin"
	db "go_blog/database"
	rd "go_blog/render"
	"strconv"
)

// get all posts
func GetPosts(c *gin.Context) {
	posts, _ := db.GetAllPostIdAndTitle()
	c.JSON(200, posts)
}

func GetPostById(c *gin.Context) {
	postIndex := c.Param("id")
	index, _ := strconv.Atoi(postIndex)
	post, _ := db.GetPostById(index)
	html := rd.RenderMd([]byte(post.Content))
	c.JSON(200, gin.H{
		"post": post,
		"html": string(html),
	})
}

func GetPrivatePostById(c *gin.Context) {
	//todo verify the user and password before get the post content
	postIndex := c.Param("id")
	index, _ := strconv.Atoi(postIndex)
	post, _ := db.GetPostById(index)
	html := rd.RenderMd([]byte(post.Content))
	c.JSON(200, gin.H{
		"post": post,
		"html": string(html),
	})
}

func AddPost(c *gin.Context) {
	// get post data
	var post db.BlogData
	// check the post data type
	if c.ContentType() != "application/json" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "post data type must be application/json",
		})
		return
	}
	c.BindJSON(&post)
	// insert post to database
	db.InsertPost(post)
	// return success
	c.JSON(200, gin.H{
		"status": "success",
	})
}
