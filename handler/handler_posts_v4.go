package handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"go_blog/database"
	renders "go_blog/render"
	"log"
	"net/http"
)

func V4Login(c *gin.Context, db_user *sql.DB) {
	type structLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type LoginResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Token   string `json:"token"`
		Name    string `json:"name"`
		Email   string `json:"email"`
	}
	var login structLogin
	if c.BindJSON(&login) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid request",
		})
		return
	}
	if database.V3Login(db_user, login.Email, login.Password) {
		c.JSON(http.StatusOK, gin.H{
			"msg":   "log in success",
			"token": database.V3GenerateToken(login.Email),
			"name":  database.GetUserByEmail(db_user, login.Email).Name,
			"email": login.Email,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "log in failed",
		})
	}
}
func V4VerifyToken(c *gin.Context, db_user *sql.DB) {
	// get auth header
	auth := c.Request.Header.Get("Authorization")
	user := database.V3GetUserByAuthHeader(db_user, auth)
	if user.Email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "invalid token",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "valid token",
		})
	}

}
func V4GetPost(c *gin.Context, db_user *sql.DB, db_post *sql.DB) {
	type GetPostRequest struct {
		Url      string `json:"url"`
		Rendered bool   `json:"rendered"`
	}
	type GetPostResponse struct {
		Status  string              `json:"status"`
		Message string              `json:"message"`
		Post    database.V4PostData `json:"post"`
		Html    string              `json:"html"`
	}
	var postRequest GetPostRequest
	var response GetPostResponse
	//if c.BindJSON(&GetPostRequest{}) != nil {
	if c.BindJSON(&postRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid request",
		})
		return
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	// get post
	postData, err := database.V4GetPostByUrlUser(db_post, postRequest.Url, user)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "permission denied",
		})
		return
	}
	response.Status = "success"
	response.Message = "ok"
	response.Post = postData
	response.Html = string(renders.RenderMd([]byte(postData.Content)))
	c.JSON(http.StatusOK, response)
}

func V4SearchPosts(c *gin.Context, db_user *sql.DB, db_post *sql.DB) {
	type SearchParams database.SearchParams
	type V4SearchPostsResponse struct {
		Status        string                `json:"status"`
		Message       string                `json:"message"`
		Posts         []database.V4PostData `json:"posts"`
		NumberOfPosts int                   `json:"number_of_posts"`
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	var searchParams SearchParams
	if c.BindJSON(&searchParams) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid request",
		})
		return
	}
	posts, err := database.V4SearchPostUser(db_post, database.SearchParams(searchParams), user)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "permission denied",
		})
		return
	}
	c.JSON(http.StatusOK, V4SearchPostsResponse{
		Status:        "success",
		Message:       "ok",
		Posts:         posts,
		NumberOfPosts: len(posts),
	})
}

func V4UpdatePost(c *gin.Context, db_user *sql.DB, db_post *sql.DB) {
	type PostUpdateRequest database.V4PostData
	type GetPostResponse struct {
		Status  string              `json:"status"`
		Message string              `json:"message"`
		Post    database.V4PostData `json:"post"`
		Html    string              `json:"html"`
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	log.Println("V4UpdatePost: user: ", user)
	var postUpdateRequest PostUpdateRequest
	if c.BindJSON(&postUpdateRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid request",
		})
		return
	}
	err := database.V4UpdatePosByUser(db_post, database.V4PostData(postUpdateRequest), user)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "permission denied",
		})
		return
	}
	postData, err := database.V4GetPostByUrlUser(db_post, postUpdateRequest.Url, user)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "permission denied",
		})
		return
	}
	c.JSON(http.StatusOK, GetPostResponse{
		Status:  "success",
		Message: "ok",
		Post:    postData,
		Html:    string(renders.RenderMd([]byte(postData.Content))),
	})
}
func V4NewPost(c *gin.Context, db_user *sql.DB, db_post *sql.DB) {
	type NewPostResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Url     string `json:"url"`
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	url, err := database.V4NewPostUser(db_post, user)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "permission denied",
		})
		return
	}
	c.JSON(http.StatusOK, NewPostResponse{
		Status:  "success",
		Message: "ok",
		Url:     url,
	})
}

func V4GetDistinct(c *gin.Context, db_user *sql.DB, db_post *sql.DB) {
	type GetDistinctRequest struct {
		Field string `json:"field"`
	}
	type GetDistinctResponse struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Values  []string `json:"values"`
		Length  int      `json:"length"`
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	var request GetDistinctRequest
	if c.BindJSON(&request) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid request",
		})
		return
	}
	print(request.Field)
	values, err := database.V4GetDistinctUser(db_post, request.Field, user)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "permission denied",
		})
		return
	}
	c.JSON(http.StatusOK, GetDistinctResponse{
		Status:  "success",
		Message: "ok",
		Values:  values,
		Length:  len(values),
	})
}
