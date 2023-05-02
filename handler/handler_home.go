package handler
import (
    "github.com/gin-gonic/gin"
    "net/http"
    db "go_blog/database"
    "strconv"
)

func HandlerHome(c *gin.Context) {
    // send index.html
    // c.HTML(http.StatusOK, "index.html", gin.H{})

    // get all posts
    posts, _ := db.GetAllPostIdAndName() 
    //tmpl, _ := template.ParseFiles("homepage.html")

    html := ""
    for index, title := range posts {
        // convert index to string int -> string Itoa
        ind := strconv.Itoa(index)
        html += `<a href="/posts/` + ind + `">` + title + `</a>`
        html += "<br>"
        _ = index
    }

    c.String(http.StatusOK, html)
}


