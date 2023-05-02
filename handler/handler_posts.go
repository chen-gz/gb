package handler
import (
    "github.com/gin-gonic/gin"
    db "go_blog/database"
    "strconv"
    rd "go_blog/render"
)
// get all posts
func HandlerGetPosts(c *gin.Context) {
    posts, _ := db.GetAllPostIdAndName() 
    c.JSON(200, posts)
}

/// This function is used to get post by id.
/// The id is the index of the post in database.
/// The index is the primary key of the post.
/// When the post is setup as private post, 
/// this function will return permission denied.
func HanlderGetPostId(c *gin.Context) {
    post_index :=  c.Param("id") 
    index, _ := strconv.Atoi(post_index)
    post, _ := db.GetPostByIndex(index)
    // todo: check if the post is private
    // if post.Private == true {
    //     c.JSON(403, gin.H{
    //         "status": "error",
    //         "message": "permission denied",
    //     })
    //     return
    // }
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
    post, _ := db.GetPostByIndex(index)
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
