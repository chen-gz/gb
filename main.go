package main

import (
	"embed"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go_blog/database"
	hd "go_blog/handler"
)

//go:embed web_src/dist/*
var frontend embed.FS

func ginServer() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Length", "Content-Type", "Origin", "Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Credentials"},
		AllowCredentials: true,
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:2009",
			"https://blog.ggeta.com"},
	}))
	///////////////////// V3 api. the database does not compatible with v2
	//r.POST("/api/v3/get_post", func(c *gin.Context) {
	//	// get url param
	//	//hd.V3GetPost(c)
	//})
	//r.POST("/api/v3/search_posts", func(c *gin.Context) {
	//	//hd.V3SearchPosts(c)
	//})
	//r.POST("/api/v3/update_post", func(c *gin.Context) {
	//	//hd.V3UpdatePost(c)
	//})
	//r.POST("/api/v3/new_post", func(c *gin.Context) {
	//	//hd.V3NewPost(c)
	//})
	//r.POST("/api/v3/get_distinct", func(c *gin.Context) {
	//	//hd.V3GetDistinct(c)
	//})
	////////////////////////// v4 api. The database is not compatible with v3
	db_blog := database.InitV4()
	db_user, _ := database.UserDbInit()
	r.POST("/api/v4/login", func(c *gin.Context) {
		//hd.V3Login(c)
		hd.V4Login(c, db_user)
	})
	r.POST("/api/v4/get_post", func(c *gin.Context) {
		hd.V4GetPost(c, db_user, db_blog)
	})

	// list all files in the frontend
	//
	//files, _ := frontend.ReadDir("web_src/dist/assets")
	//fmt.Println("files in frontend ****************************************")
	//for _, file := range files {
	//	fmt.Println(file.Name())
	//
	//}
	//
	//r.GET("/assets/*filepath", func(c *gin.Context) {
	//	//c.FileFromFS("/assets/", frontendBox)
	//	if data, err := frontend.ReadFile("web_src/dist/assets" + c.Param("filepath")); err == nil {
	//		if c.Param("filepath")[len(c.Param("filepath"))-3:] == ".js" {
	//			c.Data(200, "application/javascript", data)
	//		} else if c.Param("filepath")[len(c.Param("filepath"))-4:] == ".css" {
	//			c.Data(200, "text/css", data)
	//		} else if c.Param("filepath")[len(c.Param("filepath"))-4:] == ".svg" {
	//			c.Data(200, "image/svg+xml", data)
	//		} else {
	//			c.Data(200, "application/octet-stream", data)
	//		}
	//	} else {
	//		c.String(404, "File not found")
	//	}
	//	_, err := frontend.ReadFile("web_src/dist/assets/" + c.Param("filepath"))
	//	print(err)
	//})
	//
	//// all other path will be redirected to index.html
	////r.GET("/", func(c *gin.Context) {
	//r.NoRoute(func(c *gin.Context) {
	//	c.FileFromFS("web_src/dist/", http.FS(frontend))
	//})
	//
	r.Run(":2009") // listen and serve on

}

func main() {
	ginServer()
}
