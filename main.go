package main

import (
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go_blog/database"
	hd "go_blog/handler"
	"log"
	"net/http"
)

func gin_server() {
	r := gin.Default()
	//r.Use(cors.Default()) // allow cross origin request
	// all all cross origin request
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

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
	r.GET("/api/v1/user_get/:url", func(c *gin.Context) {
		token := c.Request.URL.Query()["token"][0]
		valid, email := hd.V1VerifyToken(token)
		if !valid {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			return
		}
		//
		url := c.Param("url")
		hd.V1UserGetPost(c, url, email)
	})
	r.POST("/api/v1/update_post/", func(c *gin.Context) {
		log.Println("update post")
		var jsonData map[string]interface{}
		if c.BindJSON(&jsonData) != nil || jsonData["post"] == nil || jsonData["token"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
			return
		}
		dataStr, err := json.Marshal(jsonData["post"])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
			return
		}
		token := jsonData["token"].(string)
		valid, email := hd.V1VerifyToken(token)
		if !valid {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			return
		}
		blogData := database.BlogDataV1{}
		if json.Unmarshal(dataStr, &blogData) != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid data structure"})
			return
		}
		hd.V1UserUpdatePost(c, blogData, email)
	})

	r.Run(":2009") // listen and serve on
}

func main() {
	gin_server()
}
