package handler
import (
    "github.com/gin-gonic/gin"
    db "go_blog/database"
    "strconv"
    rd "go_blog/render"
)
// get all posts
func HandlerGetPosts(c *gin.Context) {
    posts, _ := db.GetAllPostIdAndTitle() 
    c.JSON(200, posts)
}

func HanlderGetPostId(c *gin.Context) {
    post_index :=  c.Param("id") 
    index, _ := strconv.Atoi(post_index)
    post, _ := db.GetPostById(index)
    html := rd.RenderMd([]byte(post.Content))
    c.JSON(200, gin.H{
        "post": post,
        "html": string(html),
    })
}

func HandlerGetPrivatePostId(c *gin.Context) {
    //todo verify the user and password before get the post content 
    post_index :=  c.Param("id") 
    index, _ := strconv.Atoi(post_index)
    post, _ := db.GetPostById(index)
    html := rd.RenderMd([]byte(post.Content))
    c.JSON(200, gin.H{
        "post": post,
        "html": string(html),
    })
}

func HandlerAddPost(c *gin.Context) {
    // get post data
    var post db.BlogData
    // check the post data type
    if c.ContentType() != "application/json" {
        c.JSON(400, gin.H{
            "status": "error",
            "message": "post data type must be application/json",
        })
        return
    }
    c.BindJSON(&post)
    // insert post to database
    db.InsertPost(post)
    // return success
    c.JSON(200, gin.H{
        "status": "success",
    })
}
