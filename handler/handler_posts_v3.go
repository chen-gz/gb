package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go_blog/database"
	renders "go_blog/render"
	"net/http"
	"time"
)

type GetPostRequestV3 struct {
	Url      string `json:"url"`
	Rendered bool   `json:"rendered"` // if true, return rendered content
}
type GetPostResponseV3 struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Post    database.PostDataV2 `json:"post"`
	Html    string              `json:"html"`
}

func V3GetPost(c *gin.Context) {
	// get by url param
	result := GetPostResponseV3{
		Status: "failed",
	}
	var jsonData map[string]interface{}

	var postRequest GetPostRequestV3
	if c.BindJSON(&jsonData) != nil || mapstructure.Decode(jsonData, &postRequest) != nil {
		result.Message = "invalid request"
		c.JSON(http.StatusBadRequest, result)
		return
	}
	user := GetUserByAuthHeader(c.Request.Header.Get("Authorization"))
	post, post_content, post_comment := database.V2GetPostByUrl(postRequest.Url)
	if user.Level < post.PrivateLevel {
		result.Message = "permission denied"
		c.JSON(http.StatusForbidden, result)
		return
	}
	result.Status = "success"
	result.Message = "ok"

	result.Post = database.PostDataV2{post, post_content, post_comment}
	result.Html = string(renders.RenderMd([]byte(post_content.Content)))
	c.JSON(http.StatusOK, result)
}

type SearchPostsRequestV3 database.V2SearchParams
type SearchPostsResponseV3 struct {
	Status        string                    `json:"status"`
	Message       string                    `json:"message"`
	Posts         []database.PostDataV2Meta `json:"posts"`
	NumberOfPosts int                       `json:"number_of_posts"`
}

func V3SearchPosts(c *gin.Context) {
	result := SearchPostsResponseV3{
		Status: "failed",
	}
	// use database.V2SearchParams to search
	var jsonData map[string]interface{}
	var searchRequest SearchPostsRequestV3
	if c.BindJSON(&jsonData) != nil || mapstructure.Decode(jsonData, &searchRequest) != nil {
		result.Message = "invalid request"
		c.JSON(http.StatusBadRequest, result)
		return
	}
	user := GetUserByAuthHeader(c.Request.Header.Get("Authorization"))
	searchRequest.PrivateLevel = user.Level

	posts, cnt := database.V2SearchPosts(database.V2SearchParams(searchRequest))
	result.Status = "success"
	result.Message = "ok"
	result.Posts = posts
	result.NumberOfPosts = cnt
	c.JSON(http.StatusOK, result)
}

type UpdatePostRequestV3 database.V2UpdateParams
type UpdatePostResponseV3 struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Url     string `json:"url"` // if the url is changed, return the new url
}

func V3UpdatePost(c *gin.Context) {
	// only author and admin can update meta and post of a post
	// other registered user can only update comment
	// if user is not registered, return http.StatusForbidden
	result := UpdatePostResponseV3{
		Status: "failed",
	}
	var jsonData map[string]interface{}
	var updateRequest UpdatePostRequestV3
	if c.BindJSON(&jsonData) != nil || mapstructure.Decode(jsonData, &updateRequest) != nil {
		result.Message = "invalid request"
		c.JSON(http.StatusBadRequest, result)
		return
	}
	user := GetUserByAuthHeader(c.Request.Header.Get("Authorization"))
	if user.Level < 1 {
		result.Message = "permission denied"
		c.JSON(http.StatusForbidden, result)
		return
	}
	if user.Role != "Admin" {
		// && (updateRequest.MetaUpdate || updateRequest.CommentUpdate) {
		result.Message = "permission denied"
		c.JSON(http.StatusForbidden, result)
		return
	}
	// update post
	database.V2UpdatePost(database.V2UpdateParams(updateRequest))
	result.Status = "success"
	result.Message = "ok"
	result.Url = updateRequest.Meta.Url
	c.JSON(http.StatusOK, result)
}

type NewPostResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Url     string `json:"url"`
}

func V3NewPost(c *gin.Context) {
	response := NewPostResponse{
		Status: "failed",
	}
	auth := c.Request.Header.Get("Authorization")
	user := GetUserByAuthHeader(auth)
	if user.Role != "admin" {
		response.Message = "permission denied"
		c.JSON(http.StatusForbidden, response)
	}
	post := database.PostDataV2{
		database.PostDataV2Meta{
			Url:        time.Now().String(),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		},
		database.PostDataV2Content{},
		database.PostDataV2Comment{},
	}

	err := database.V2InsertPost(post)
	if err != nil {
		response.Message = "internal error"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Status = "success"
	response.Message = "ok"
	response.Url = post.Meta.Url
}

type GetDistinctRequest struct {
	Column string `json:"column"`
}
type GetDistinctResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Values  []string `json:"values"`
	Length  int      `json:"length"`
}

func V3GetDistinct(c *gin.Context) {
	response := GetDistinctResponse{
		Status: "failed",
	}
	var jsonData map[string]interface{}
	var request GetDistinctRequest
	if c.BindJSON(&jsonData) != nil || mapstructure.Decode(jsonData, &request) != nil {
		response.Message = "invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	values, err := database.V2GetDistinct(request.Column)
	if err != nil {
		response.Message = "internal error"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Status = "success"
	response.Message = "ok"
	response.Values = values
	response.Length = len(values)
	c.JSON(http.StatusOK, response)
}
