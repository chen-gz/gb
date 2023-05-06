package handler

import (
	"github.com/gin-gonic/gin"
	db "go_blog/database"
	rd "go_blog/render"
	"math"
	"net/http"
)

func V1GetPost(c *gin.Context, url string) {
	post := db.V1GetPostByUrl(url)
	post.Rendered = string(rd.RenderMd([]byte(post.Content)))
	c.JSON(http.StatusOK, post)
}

func V1SearchPosts(c *gin.Context) {
	// get query params as map[string][]string
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		params[k] = v[0]
	}
	summary := params["summary"] == "true"
	posts := db.V1SearchPost(params)
	if summary {
		for i := 0; i < len(posts); i++ {
			posts[i].Summary = posts[i].Content[:int(math.Min(float64(len(posts[i].Content)), 500))]
			posts[i].Rendered = string(rd.RenderMd([]byte(posts[i].Summary)))
			posts[i].Content = ""
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}
