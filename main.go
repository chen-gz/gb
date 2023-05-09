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
		// get token and verify token
		token := c.Request.URL.Query()["token"][0]
		valid, _ := hd.V1VerifyToken(token)

		if !valid {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "invalid token",
			})
			return
		}
		//
		url := c.Param("url")
		hd.V1UserGetPost(c, url, token)
	})
	r.POST("/api/v1/update_post/", func(c *gin.Context) {
		log.Println("update post")

		var jsonData map[string]interface{}
		c.BindJSON(&jsonData)
		data := jsonData["post"]
		//OOO := data.(map[string]interface{})
		token := jsonData["token"].(string)
		valid, email := hd.V1VerifyToken(token)
		log.Println("valid: ", valid)
		log.Println("email: ", email)
		blogData := database.BlogDataV1{}
		//blogData = OOO.(database.BlogDataV1)
		// convert map[string]interface{} to string
		dataStr, _ := json.Marshal(data)
		json.Unmarshal(dataStr, &blogData)
		//json.Unmarshal([]byte(data.(string)), &blogData)

		log.Println("blogData: ", blogData)

		// convert OOO to blogData
		//blogData = OOO as BlogDataV1
		//if OOO["Id"] != nil {
		//	blogData.Id = int(OOO["Id"].(float64))
		//}
		//if OOO["Author"] != nil {
		//	blogData.Author = OOO["Author"].(string)
		//}
		//if OOO["Title"] != nil {
		//	blogData.Title = OOO["Title"].(string)
		//}
		//if OOO["Content"] != nil {
		//	blogData.Content = OOO["Content"].(string)
		//}
		//if OOO["Tags"] != nil {
		//	blogData.Tags = OOO["Tags"].(string)
		//}
		//if OOO["Categories"] != nil {
		//	blogData.Categories = OOO["Category"].(string)
		//}
		//if OOO["Url"] != nil {
		//	blogData.Url = OOO["Url"].(string)
		//}
		//if OOO["Like"] != nil {
		//	blogData.Like = int(OOO["Like"].(float64))
		//}
		//if OOO["Dislike"] != nil {
		//	blogData.Dislike = int(OOO["Dislike"].(float64))
		//}
		//if OOO["CoverImg"] != nil {
		//	blogData.CoverImg = OOO["CoverImg"].(string)
		//}
		//if OOO["IsDraft"] != nil {
		//	blogData.IsDraft = OOO["IsDraft"].(bool)
		//}
		//if OOO["IsDeleted"] != nil {
		//	blogData.IsDeleted = OOO["IsDeleted"].(bool)
		//}
		//if OOO["PrivateLevel"] != nil {
		//	blogData.PrivateLevel = int(OOO["PrivateLevel"].(float64))
		//}
		//if OOO["ViewCount"] != nil {
		//	blogData.ViewCount = int(OOO["ViewCount"].(float64))
		//}
		//if OOO["CreatedAt"] != nil {
		//	str := OOO["CreatedAt"].(string)
		//	// convert string (str) to time.Time
		//	t, _ := time.Parse("2006-01-02 15:04:05", str)
		//	blogData.CreatedAt = t
		//}
		//if OOO["UpdatedAt"] != nil {
		//	str := OOO["UpdatedAt"].(string)
		//	// convert string (str) to time.Time
		//	t, _ := time.Parse("2006-01-02 15:04:05", str)
		//	blogData.UpdatedAt = t
		//}

		//log.Println("blogData: ", blogData)

		hd.V1UserUpdatePost(c, blogData, token)
		if !valid {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "invalid token",
			})
		}
	})

	r.Run(":2009") // listen and serve on
}

func main() {
	gin_server()
}
