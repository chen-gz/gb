package handler
import (
    "github.com/gin-gonic/gin"
    "net/http"
)
func HandlerAdmin(c *gin.Context) {
    c.String(http.StatusOK, "")
}
func HandlerAdminLogin(c *gin.Context) {
}
func HandlerAdminLogout(c *gin.Context) {
}
func HandlerAdminAddPost(c *gin.Context) {
}
func HandlerAdminEditPost(c *gin.Context) {
}
func HandlerAdminDeletePost(c *gin.Context) {
}
func HandlerAdminGetPosts(c *gin.Context) {
}
func HandlerAdminAddPage(c *gin.Context) {
}

