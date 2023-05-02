package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    hd "go_blog/handler"
)

func  gin_server() {
    r := gin.Default()
    r.Use(cors.Default())  // allow cross origin request
    r.GET("/api/", func(c *gin.Context) {
        hd.HandlerHome(c)
    })
    r.GET("/api/admin", func(c *gin.Context) {
        hd.HandlerAdmin(c)
    })
    r.GET("/api/posts", func(c *gin.Context) {
        hd.HandlerGetPosts(c)
    })
    r.GET("/api/posts/:id", func(c *gin.Context) {
        hd.HanlderGetPostId(c)
    })
    r.POST("/api/posts", func(c *gin.Context) {
        hd.HandlerAddPost(c)
    })
    // uplaod image and file to server and return a link to client 
    r.POST("/api/upload", func(c *gin.Context) {
        hd.HandlerUpload(c)
    })
    r.Run() // listen and serve on
}


func main() {
    gin_server()
}
