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
	r.POST("/api/v2/search_posts", func(c *gin.Context) {
		// parameter rqueset,
		var searchRequest map[string]interface{}
		//token := c.Request.URL.Query().Get("token")
		//searchRequest =
		//token
		c.BindJSON(&searchRequest)
		log.Println("searchRequest", searchRequest)
		token := searchRequest["token"].(string)

		//token = c.Request.Pa
		var search_params database.SearchParams
		//err := json.Unmarshal([]byte(searchRequest["search_params"].(string)), &search_params)
		dataStr, err := json.Marshal(searchRequest["search_params"])
		if err != nil {
			log.Println("error when marshal search_params")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "error when marshal search_params",
			})
		}
		if json.Unmarshal(dataStr, &search_params) != nil {
			log.Println("error when unmarshal search_params")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "error when unmarshal search_params",
			})
		}

		//verify token
		if token != "" {

			valid, email := hd.V1VerifyToken(token)
			_ = email
			if !valid {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "invalid token",
				})
				return
			}
		}
		result := database.V1SearchPostBySearchParams(search_params)
		c.JSON(http.StatusOK, result)

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
	r.GET("/api/v1/user_verify", func(c *gin.Context) {
		hd.V1Verify(c)
	})

	r.Run(":2009") // listen and serve on
}

func main() {
	gin_server()
}
