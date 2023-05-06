package handler

import (
	"github.com/gin-gonic/gin"
	db "go_blog/database"
)

func V1GetCategories(c *gin.Context) {
	cates := db.V1GetCategories()
	c.JSON(200, cates)
}
