package handler

import (
	"github.com/gin-gonic/gin"
	db "go_blog/database"
	rd "go_blog/render"
	"net/http"
)

//	type BlogData struct {
//		Id         int
//		Author     string
//		Title      string
//		Content    string
//		Tags       []string
//		Categories []string
//		Datetime   time.Time
//		Url        string // for vue router and s3 storage. no space.
//	}
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func Home(c *gin.Context) {
	// get 20 recent posts
	posts, _ := db.GetRecentPosts(20)
	retPosts := make([]map[string]interface{}, 0)
	for i := 0; i < len(posts); i++ {
		retPost := make(map[string]interface{})
		retPost["Id"] = posts[i].Id
		retPost["Author"] = posts[i].Author
		retPost["Title"] = posts[i].Title
		retPost["Content"] = posts[i].Content
		retPost["Tags"] = posts[i].Tags
		retPost["Categories"] = posts[i].Categories
		retPost["Datetime"] = posts[i].Datetime
		retPost["Url"] = posts[i].Url
		retPost["Html"] = string(rd.RenderMd([]byte(posts[i].Content)))
		retPost["Summary"] = string(rd.RenderMd([]byte(posts[i].Content[:min(500, len(posts[i].Content))])))
		retPosts = append(retPosts, retPost)

	}
	c.JSON(http.StatusOK, retPosts)
}
