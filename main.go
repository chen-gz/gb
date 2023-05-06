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

	r.Run(":2009") // listen and serve on
}

func main() {
	gin_server()
}
