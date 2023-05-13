package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	db "go_blog/database"
	rd "go_blog/render"
	"io/ioutil"
	"log"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid request",
			"post":    db.BlogDataV1{}})
		return
	}
	auth := c.Request.Header.Get("Authorization")
	user := GetUserByAuthHeader(auth)

	post := db.V1GetPostByUrl(postRequest.Url)
	// todo: introduce group to post || user.Group == post.Group {
	if user.Group == "admin" || user.Level >= post.PrivateLevel {

	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"message": "permission denied",
			"post":    db.BlogDataV1{}})
		return
	}
	if postRequest.Summary {
		post.Summary = post.Content[:int(math.Min(float64(len(post.Content)), 500))]
	}
	if postRequest.Rendered {
		if postRequest.Summary {
			post.Rendered = string(rd.RenderMd([]byte(post.Summary)))
		} else {
			post.Rendered = string(rd.RenderMd([]byte(post.Content)))
		}
	}
	if !postRequest.Content {
		post.Content = ""
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "post found",
		"post":    post,
	})
}

type NewPostRequest struct {
	Token string `json:"token"`
}

func V2NewPost(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	user := GetUserByAuthHeader(auth)
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"message": "permission denied",
			"post":    db.BlogDataV1{},
		})
		return
	}
	post := db.BlogDataV1{
		Title:     "New Post",
		Url:       time.Now().String(), // url should be unique
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// if this get error, generate new url and try again
	tries := 0
	for db.V1InsertPost(post) != nil {
		tries++
		post.Url = time.Now().String()
		if tries > newPostRetry {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "internal error, please try again later",
				"post":    db.BlogDataV1{},
			})
			return
		}
	}
	post = db.V1GetPostByUrl(post.Url)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "new post created",
		"post":    post},
	)
}

type DeletePostRequest struct {
	Url string `json:"url"`
}

func V2DeletePost(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	user := GetUserByAuthHeader(auth)
	log.Println("user :", user)

	var jsonData map[string]interface{}
	var postRequest DeletePostRequest
	if c.BindJSON(&jsonData) != nil || mapstructure.Decode(jsonData, &postRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status:": "failed",
			"message": "invalid data"})
		return
	}
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"message": "permission denied"})
		return
	}
	// delete post
	if db.V1DeletePost(postRequest.Url) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "post deleted"})
}

type RenderMdRequest struct {
	Token   string `json:"token"`
	Content string `json:"content"`
}

func V2SearchPost(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	user := GetUserByAuthHeader(auth)
	_ = user
	var searchRequest map[string]interface{}
	c.BindJSON(&searchRequest)
	var search_params db.SearchParams
	dataStr, err := json.Marshal(searchRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "error when marshal search_params",
			"posts":   []db.BlogDataV1{},
		})
		return
	}
	if json.Unmarshal(dataStr, &search_params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "error when unmarshal search_params",
			"posts":   []db.BlogDataV1{},
		})
		return
	}
	// todo: some result are not allow to see for all user
	result, lengh := db.V1SearchPostBySearchParams(search_params)
	if search_params.Summary == true && search_params.Rendered == true {
		for i := range result {
			result[i].Rendered = string(rd.RenderMd([]byte(result[i].Summary)))
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "search success",
		"posts":   result,
		"count":   lengh,
	})

}
func V2UpdatePost(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	user := GetUserByAuthHeader(auth)
	var jsonData map[string]interface{}
	var postRequest db.BlogDataV1
	if c.BindJSON(&jsonData) != nil || mapstructure.Decode(jsonData, &postRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid data",
		})
		return
	}
	// check permission
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"message": "permission denied",
		})
		return
	}
	// update post
	postRequest.UpdatedAt = time.Now()
	if db.V1UpdatePost(postRequest) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "internal error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "post updated",
	})
}

func V2RenderMd(c *gin.Context) {
	// get auth header
	auth := c.Request.Header.Get("Authorization")
	user := GetUserByAuthHeader(auth)

	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "permission denied",
			"html":   "",
		})
		return
	}
	// get body to string
	body, _ := ioutil.ReadAll(c.Request.Body)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"html":   string(rd.RenderMd(body)),
	})
}
