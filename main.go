package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	hd "go_blog/handler"
	"log"
)

func gin_server() {
	r := gin.Default()
	r.Use(cors.Default()) // allow cross origin request
	// all all cross origin request
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

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
	r.POST("/api/v1/update_post/", func(c *gin.Context) {
		log.Println("update post")
		//data := database.BlogDataV1{}
		//_ = c.BindJSON(&data)

		// get auth token from header
		token := c.GetHeader("Authorization")
		log.Println(c.Request.Header)
		log.Println(c.Request.Body)
		// get json data from body
		var json map[string]interface{}
		c.BindJSON(&json)
		log.Println(json)

		log.Println(c.Get("firstName"))
		log.Println("token: ", token)
		//hd.V1UserUpdatePost(c, data, token)
	})

	r.Run(":2009") // listen and serve on
}

func main() {
	gin_server()
}
