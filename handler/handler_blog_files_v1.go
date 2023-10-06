package handler

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go_blog/database"
	"net/http"
	"path/filepath"
	"time"
)

const endpoint = "minio.ggeta.com"
const accessKeyID = "HI4mSQabJ6GWesqES4V4"
const secreteAccessKey = "WIK6SwKqceiPCalmhDj4meOdqLdErSfw4QNpEZxx"
const bucketName = "blog-public-data"

func GetPresignedUrl(c *gin.Context, db_user *sql.DB) {
	type UploadFileRequest struct {
		FileName  string `json:"file_name"`
		HashCrc32 string `json:"hash_crc32"`
	}
	type UploadFileResponse struct {
		PresignedUrl string `json:"presigned_url"`
		Message      string `json:"message"`
		Filename     string `json:"filename"` // the file name with be update by the server
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	if user.Email == "" {
		c.JSON(http.StatusForbidden, UploadFileResponse{
			Message: "permission denied",
		})
		return
	}
	// ok get request
	var uploadFileRequest UploadFileRequest
	if c.BindJSON(&uploadFileRequest) != nil {
		c.JSON(http.StatusBadRequest, UploadFileResponse{
			Message: "invalid request",
		})
		return
	}
	// file_name append with hashCrc32
	//file_name_with_hash := uploadFileRequest.FileName + "_" + uploadFileRequest.HashCrc32 // this should not exist in minio (skip check here)
	// sperate file name and file extension

	filename := uploadFileRequest.FileName
	extension := filepath.Ext(filename)
	nameWithoutExtension := filename[:len(filename)-len(extension)]
	file_name_with_hash := nameWithoutExtension + "_" + uploadFileRequest.HashCrc32 + filepath.Ext(filename)

	client, err := MinioClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	presignedURL, err := client.PresignedPutObject(c, bucketName, file_name_with_hash, time.Hour) // 1 hour expiry
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Presigned URL for uploading: ", presignedURL)
	c.JSON(http.StatusOK, UploadFileResponse{
		PresignedUrl: presignedURL.String(),
		Filename:     file_name_with_hash,
		Message:      "success",
	})
}
func MinioClient() (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secreteAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return minioClient, nil
}
