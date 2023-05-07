package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	hd "go_blog/handler"
)

func gin_server() {
	r := gin.Default()
	r.Use(cors.Default()) // allow cross origin request

	r.GET("/api/v1/get_post/:url", func(c *gin.Context) {
		hd.V1GetPost(c, c.Param("url"))
	})
	r.GET("/api/v1/search_posts", func(c *gin.Context) {
		hd.V1SearchPosts(c)
	})
	r.GET("/api/v1/get_tags", func(c *gin.Context) {
		hd.V1GetTags(c)
	})
	r.GET("/api/v1/get_categories", func(c *gin.Context) {
		hd.V1GetCategories(c)
	})
	r.GET("/api/v1/login", func(c *gin.Context) {
		hd.V1Login(c)
	})

	r.Run(":2009") // listen and serve on
}

func main() {
	gin_server()
}
