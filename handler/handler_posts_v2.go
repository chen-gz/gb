package handler

//
//import (
//	"github.com/gin-gonic/gin"
//	rd "go_blog/render"
//	"io/ioutil"
//	"net/http"
//)
//
//type DeletePostRequest struct {
//	Url string `json:"url"`
//}
//
//type RenderMdRequest struct {
//	Token   string `json:"token"`
//	Content string `json:"content"`
//}
//
//func V2RenderMd(c *gin.Context) {
//	// get auth header
//	auth := c.Request.Header.Get("Authorization")
//	user := GetUserByAuthHeader(auth)
//
//	if user.Role != "admin" {
//		c.JSON(http.StatusForbidden, gin.H{
//			"status": "permission denied",
//			"html":   "",
//		})
//		return
//	}
//	// get body to string
//	body, _ := ioutil.ReadAll(c.Request.Body)
//	c.JSON(http.StatusOK, gin.H{
//		"status": "success",
//		"html":   string(rd.RenderMd(body)),
//	})
//}
