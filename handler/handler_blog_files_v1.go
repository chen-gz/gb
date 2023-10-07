package handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go_blog/database"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

//const endpoint = "minio.ggeta.com"
//const accessKeyID = "HI4mSQabJ6GWesqES4V4"
//const secreteAccessKey = "WIK6SwKqceiPCalmhDj4meOdqLdErSfw4QNpEZxx"
//const bucketName = "blog-public-data"

type MinioConfig struct {
	Endpoint         string `json:"endpoint"`
	AccessKeyID      string `json:"access_key_id"`
	SecreteAccessKey string `json:"secrete_access_key"`
	BucketName       string `json:"bucket_name"`
}

var config MinioConfig

func GetPresignedUrl(c *gin.Context, db_user *sql.DB, db_blog *sql.DB, client *minio.Client) {
	type UploadFileRequest struct {
		FileName  string `json:"file_name"`
		PostId    int    `json:"post_id"`
		HashCrc32 string `json:"hash_crc32"`
	}
	type UploadFileResponse struct {
		PresignedUrl string `json:"presigned_url"`
		Message      string `json:"message"`
		Filename     string `json:"filename"` // the file name with be update by the server
		FileUrl      string `json:"file_url"` // the file url with be update by the server
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	var uploadFileRequest UploadFileRequest
	if c.BindJSON(&uploadFileRequest) != nil {
		c.JSON(http.StatusBadRequest, UploadFileResponse{Message: "invalid request"})
		return
	}
	if !database.UpdatePostPermissionCheck(db_blog, user, uploadFileRequest.PostId) {
		c.JSON(http.StatusForbidden, UploadFileResponse{Message: "permission denied"})
		return
	}

	filename := uploadFileRequest.FileName
	extension := filepath.Ext(filename)
	nameWithoutExtension := filename[:len(filename)-len(extension)]
	file_name_with_hash := nameWithoutExtension + "_" + uploadFileRequest.HashCrc32 + filepath.Ext(filename)
	publicUrl := config.Endpoint + "/" + config.BucketName + "/" + file_name_with_hash
	err := database.InsertFileUser(db_blog, user, uploadFileRequest.PostId, file_name_with_hash, publicUrl)

	presignedURL, err := client.PresignedPutObject(c, config.BucketName, file_name_with_hash, time.Hour) // 1 hour expiry
	if err != nil {
		log.Println("GetPresignedUrl: ", err)
		return
	}
	log.Println("Presigned URL for uploading: ", presignedURL)

	c.JSON(http.StatusOK, UploadFileResponse{
		PresignedUrl: presignedURL.String(),
		Filename:     file_name_with_hash,
		Message:      "success",
		FileUrl:      publicUrl,
	})
}

func InitMinioClient(_config MinioConfig) *minio.Client {
	config = _config
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecreteAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatal("InitMinioClient: ", err)
		return nil
	}
	return minioClient
}
