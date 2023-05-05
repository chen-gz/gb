package handler

import (
	"github.com/gin-gonic/gin"
	db "go_blog/database"
	rd "go_blog/render"
	"strconv"
)

// get all posts
func GetPosts(c *gin.Context) {
	posts, _ := db.GetAllPostIdAndTitle()
	c.JSON(200, posts)
}

func GetPostById(c *gin.Context) {
	postIndex := c.Param("id")
	index, _ := strconv.Atoi(postIndex)
	post, _ := db.GetPostById(index)
	html := rd.RenderMd([]byte(post.Content))
	c.JSON(200, gin.H{
		"post": post,
		"html": string(html),
	})
}

func GetPrivatePostById(c *gin.Context) {
	//todo verify the user and password before get the post content
	postIndex := c.Param("id")
	index, _ := strconv.Atoi(postIndex)
	post, _ := db.GetPostById(index)
	html := rd.RenderMd([]byte(post.Content))
	c.JSON(200, gin.H{
		"post": post,
		"html": string(html),
	})
}

func AddPost(c *gin.Context) {
	// get post data
	var post db.BlogData
	// check the post data type
	if c.ContentType() != "application/json" {
		c.JSON(400, gin.H{
			"status":  "error",
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

// id, title, content, created_at, updated_at, 
// is_private, is_draft, is_deleted, 
// tags, category, author, comments, likes, 
// views, cover_image


func V1GetPost(c *gin.Context, url string) {
    url := c.Param("url")
	database, err := sql.Open(dbType, dbPath)
	rows, err := database.Query("SELECT * FROM posts WHERE url = ?", url)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "error",
            "message": "post not found",
        })
        return
    }
    defer rows.Close()
    post := BlogData{}

	tag := ""
	category := ""
	err = query.Scan(&post.Id, &post.Title, &post.Author,
		&post.Content, &tag, &category,
		&post.Datetime, &post.Url)
	post.Tags = strings.Split(tag, ",")
	post.Categories = strings.Split(category, ",")

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "error",
            "message": "post not found",
        })
        return
    }

    html := rd.RenderMd([]byte(post.Content))
    c.JSON(200, gin.H{
        "title": post.Title,
        "author": post.Author,
        "datetime": post.Datetime,
        "tags": post.Tags,
        "categories": post.Categories,
        "content": post.Content,
        "html": string(html),
    })
}


// /api/v1/get_post_by_id/:id?(params)
//      if post is private or not exist, return 404
//      if post is public, return post
//
//           :limit <-- (int) limit the number of posts
//           :content <-- (bool) return content or not
//           ... other post fields
