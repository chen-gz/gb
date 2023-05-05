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
        hd.V1GetPost(c, c.Param("url"));
    })
    r.GET("/api/v1/search_posts/", func(c *gin.Context) {
        hd.V1SearchPosts(c);
    })
	r.GET("/api/home", func(c *gin.Context) {
		hd.Home(c)
	})
	r.GET("/api/admin", func(c *gin.Context) {
		hd.HandlerAdmin(c)
	})
	r.GET("/api/posts", func(c *gin.Context) {
		hd.GetPosts(c)
	})
	r.GET("/api/posts/:id", func(c *gin.Context) {
		hd.GetPostById(c)
	})
	r.POST("/api/posts", func(c *gin.Context) {
		hd.AddPost(c)
	})
	// uplaod image and file to server and return a link to client
	r.POST("/api/upload", func(c *gin.Context) {
		hd.HandlerUpload(c)
	})





	r.Run(":2009") // listen and serve on
}

func main() {
	gin_server()
}
// a post contains:
// id, title, content, created_at, updated_at, 
// is_private, is_draft, is_deleted, 
// tags, category, author, comments, likes, 
// views, cover_image, images, files

// version 1 support api
// GET 

// /api/v1/get_post_by_id/:id?(params)
//      if post is private or not exist, return 404
//      if post is public, return post
//
//           :limit <-- (int) limit the number of posts
//           :content <-- (bool) return content or not
//           ... other post fields
//
// /api/v1/get_post_by_tag/:tag?(params)
//      if tag is not exist, return 404
//      if tag is exist, return all public posts
//
//      if params can be follwing things:
//           :limit <-- (int) limit the number of posts
//           :content <-- (bool) return content or not
//           ... other post fields
//      if tag is not exist, return 404
//











