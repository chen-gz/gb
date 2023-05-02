package main

import (
    rd "go_blog/render"
    _ "fmt"                        // no more error
    "github.com/gin-gonic/gin"
    _ "net/http"
    _ "path/filepath"
    _ "time"
    "log"
    _ "strings"
    _ "io/ioutil"
    "text/template"
    _ "os"
    db "go_blog/database"
    "strconv"
    "github.com/gin-contrib/cors"
)



func handler_homepage(c *gin.Context) {
    // send index.html
    // c.HTML(http.StatusOK, "index.html", gin.H{})

    // get all posts
    posts, _ := db.GetAllPostIdAndName() 
    tmpl, _ := template.ParseFiles("homepage.html")

    html := ""
    for index, title := range posts {
        // convert index to string int -> string Itoa
        ind := strconv.Itoa(index)
        html += `<a href="/posts/` + ind + `">` + title + `</a>`
        html += "<br>"
        _ = index
    }

    tmpl.Execute(c.Writer, html)
    // make an empty page
    //c.String(http.StatusOK, "")
    //c.String(http.StatusOK, html)
}


func handler_posts(c *gin.Context) {
    post_index :=  c.Param("id") // post_name should be an unique string
    // convert string to int
    log.Println(post_index)
    index, _ := strconv.Atoi(post_index)
    log.Println("render post: ", index)
    // get post by index
    post, _ := db.GetPostByIndex(index)
    log.Println(post.Title)
    // log.Println(post.Content)
    // convert content from string to []byte
    content := rd.RenderMd([]byte(post.Content))
    tmpl, _ := template.ParseFiles("posts.html")
    BlogData := db.BlogData{Title: post.Title, Content: string(content)}
    tmpl.Execute(c.Writer, BlogData)
}
func handler_admin(c *gin.Context) {
    // send admin.html
    // c.HTML(http.StatusOK, "admin.html", gin.H{})
    tmpl, _ := template.ParseFiles("admin.html")
    tmpl.Execute(c.Writer, nil)
}

func  gin_server() {
    r := gin.Default()
    r.Use(cors.Default())
    r.GET("/", func(c *gin.Context) {
        // gin server for home page
        handler_homepage(c)
        log.Println("home page")
    })

    r.GET("/posts/:id", func(c *gin.Context) {
        handler_posts(c)
    })
    r.GET("/admin", func(c *gin.Context) {
        handler_admin(c)
    })

    r.Run() // listen and serve on
}


func main() {
    gin_server()
    // db.Init()
    // get test data and insert to database
    // var blog = db.BlogData{5, "test232sfdssfd3423", "tessfsdft", "test", []string{"test"}, []string{"test"}, time.Now()}
    // db.InsertPost(blog)
}
