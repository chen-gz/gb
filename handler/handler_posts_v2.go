package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	db "go_blog/database"
	rd "go_blog/render"
	"io/ioutil"
	"log"
	"net/http"
)

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
