package main

//import (
//"log"
//"time"
// "fmt"
//"io/ioutil"
//"github.com/minio/minio-go/v7"
//"github.com/minio/minio-go/v7/pkg/credentials"
//"context"
// "github.com/gomarkdown/markdown"
// auth "github.com/abbot/go-http-auth"
//)
// import "github.com/gin-gonic/gin"
// import "net/http"
// import local package "render"
import (
    rd "go_blog/render"
    "fmt"
)


// func render_file() {
//     // read file
//     file, _ := ioutil.ReadFile("test.md")
//     // make sure the connect between "$" and "$$" is not change
// 
//     // render use gomarkdown
//     html := markdown.ToHTML(file, nil, nil)
//     // append mathjax script to html
//     html = append(html, []byte(`<script type="text/javascript" async src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.5/MathJax.js?config=TeX-MML-AM_CHTML"></script>`)...)
//     // write file
//     ioutil.WriteFile("test.html", html, 0644)
// }



// gin
// use minio to store the file
// use gomarkdown to render the markdown file
// use mathjax to render the math formula
// use sqlite to store the blog data
// use gin to build the web server

// homepace handler
func handler_homepage(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{})

    html, _ := rd.RenderMd("test.md")
    fmt.Println(string(html))
}

func handler_posts(c *gin.Context) {

	post_name := c.Param("id")
    // get post file
    // render post 
    // return html
}

func  gin_server() {
    // gin server for home page
    r := gin.Default()
    // r.LoadHTMLGlob("index.html")
    r.GET("/", func(c *gin.Context) {
        handler_homepage(c)
    })

	r.GET("/users/:id", func(c *gin.Context) {
		c.String(200, "The user id is  %s", id)
	})

    r.GET("/post/:post_name", func(c *gin.Context) {
        html := handler_posts(c, post_name)
        c.HTML(http.StatusOK, "index.html", gin.H{})

    })
    r.Run() // listen and serve on
}

func main() {
    // html, _ := rd.RenderMd("test.md")
    // fmt.Println(string(html))


    // fmt.Println("done")
    gin_server()
}
