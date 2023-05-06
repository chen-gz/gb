package handler

import (
	"github.com/gin-gonic/gin"
	db "go_blog/database"
	"net/http"
)

func V1GetTags(c *gin.Context) {
	tags := db.V1GetTags()
	c.JSON(http.StatusOK, tags)

}
func HandlerAddTag(c *gin.Context) {
}
func HandlerEditTag(c *gin.Context) {
}
