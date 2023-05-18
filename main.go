package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	hd "go_blog/handler"
)

func gin_server() {
	r := gin.Default()
	//r.Use(cors.Default()) // allow cross origin request
	// all all cross origin request
	r.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Length", "Content-Type", "Origin", "Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Credentials"},

		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3000", "http://blog.ggeta.com", "https://blog.ggeta.com"},
	}))
	r.GET("/api/v1/get_tags", func(c *gin.Context) {
		hd.V1GetTags(c)
	})
	r.GET("/api/v1/get_categories", func(c *gin.Context) {
		hd.V1GetCategories(c)
	})
	r.POST("/api/v2/login", func(c *gin.Context) {
		hd.V2Login(c)
	})
	// bellow are finished api
	r.POST("/api/v2/update_post", func(c *gin.Context) {
		hd.V2UpdatePost(c)
	})
	r.POST("/api/v2/search_posts", func(c *gin.Context) {
		hd.V2SearchPost(c)
	})
	r.POST("/api/v2/get_post", func(c *gin.Context) {
		hd.V2GetPost(c)
	})
	r.POST("/api/v2/delete_post", func(c *gin.Context) {
		hd.V2DeletePost(c)
	})
	r.POST("/api/v2/new_post", func(c *gin.Context) {
		hd.V2NewPost(c)
	})
	r.POST("/api/v2/render_md", func(c *gin.Context) {
		hd.V2RenderMd(c)
	})
	///////////////////// V3 api. the database does not compatible with v2
	///////////////////// so v2 will be deprecated
	r.POST("/api/v3/get_post", func(c *gin.Context) {
		// get url param
		hd.V3GetPost(c)

	})
	r.POST("/api/v3/search_posts", func(c *gin.Context) {

	})
	r.POST("/api/v3/update_post", func(c *gin.Context) {

	})
	r.POST("/api/v3/new_post", func(c *gin.Context) {

	})

	r.Run(":2009") // listen and serve on
}

func main() {
	gin_server()
}
