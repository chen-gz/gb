package handler
import (
    "github.com/gin-gonic/gin"
)

func HandlerUpload(c *gin.Context) {
    // get file
    file, _ := c.FormFile("file")
    // save file to server
    c.SaveUploadedFile(file, "static/" + file.Filename)
    // return a link to client
    c.JSON(200, gin.H{
        "status": "success",
        "link": "http://localhost:8080/static/" + file.Filename,
    })
}
