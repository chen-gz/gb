package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	db "go_blog/database"
	rd "go_blog/render"
	"math"
	"net/http"
	"time"
)

const newPostRetry = 3

func getUserFromToken(token string) (db.UserData, bool) {
	if token == "" {
		return db.UserData{}, true
	}
	valid, email := V1VerifyToken(token)
	if !valid {
		return db.UserData{}, false
	}
	return db.GetUser(email), true
}

// content, summary, rendered only return when request.
// This is in order to save bandwidth.
type GetPostRequest struct {
	Token    string `json:"token"`
	Url      string `json:"url"`
	Content  bool   `json:"content"`  // if true, return content
	Summary  bool   `json:"summary"`  // if true, return summary
	Rendered bool   `json:"rendered"` // if true, return rendered content
}

// this is handler of /api/v2/get_post with POST method
func V2GetPost(c *gin.Context) {
	var jsonData map[string]interface{}
	var postRequest GetPostRequest
	if c.BindJSON(&jsonData) != nil || mapstructure.Decode(jsonData, &postRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	// get user based on token
	user := db.GetUser("")
	if postRequest.Token != "" {
		valid, email := V1VerifyToken(postRequest.Token)
		if !valid {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			return
		}
		user = db.GetUser(email)
	}
	// get post based on user, Role , PrivateLevel and Group
	post := db.V1GetPostByUrl(postRequest.Url)
	// todo: introduce group to post || user.Group == post.Group {
	if user.Group == "admin" || user.Level >= post.PrivateLevel {

	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}
	if postRequest.Summary {
		post.Summary = post.Content[:int(math.Min(float64(len(post.Content)), 500))]
	}
	if postRequest.Rendered {
		post.Rendered = string(rd.RenderMd([]byte(post.Content)))
	}
	if !postRequest.Content {
		post.Content = ""
	}
	c.JSON(http.StatusOK, gin.H{
		"post":   post,
		"status": "success",
	})
}

type NewPostRequest struct {
	Token string `json:"token"`
}

func V2NewPost(c *gin.Context) {
	var jsonData map[string]interface{}
	var postRequest NewPostRequest
	if c.BindJSON(&jsonData) != nil || mapstructure.Decode(jsonData, &postRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	// get user based on token
	user, valid := getUserFromToken(postRequest.Token)
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		return
	}
	// check user permission
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}
	// generate url
	post := db.BlogDataV1{
		Url: time.Now().String(),
	}
	// if this get error, generate new url and try again
	tries := 0
	for db.V1InsertPost(post) != nil {
		tries++
		post.Url = time.Now().String()
		if tries > newPostRetry {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  "failed to generate url",
			})
			return
		}
	}
	post = db.V1GetPostByUrl(post.Url)
	c.JSON(http.StatusOK, gin.H{"status": "success", "post": post})
}

type DeletePostRequest struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

func V2DeletePost(c *gin.Context) {
	var jsonData map[string]interface{}
	var postRequest DeletePostRequest
	if c.BindJSON(&jsonData) != nil || mapstructure.Decode(jsonData, &postRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	// get user based on token
	user, valid := getUserFromToken(postRequest.Token)
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		return
	}
	// check user permission
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}
	// delete post
	if db.V1DeletePost(postRequest.Url) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

type RenderMdRequest struct {
	Token   string `json:"token"`
	Content string `json:"content"`
}

func V2Render(c *gin.Context) {
	//verify token is valid and user is admin
	var jsonData map[string]interface{}
	var renderRequest RenderMdRequest
	if c.BindJSON(&jsonData) != nil || mapstructure.Decode(jsonData, &renderRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	// get user based on token
	user, valid := getUserFromToken(renderRequest.Token)
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		return
	}
	// check user permission

	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"html":   string(rd.RenderMd([]byte(renderRequest.Content))),
	})
}
