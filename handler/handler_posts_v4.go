package handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"go_blog/database"
	renders "go_blog/render"
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
