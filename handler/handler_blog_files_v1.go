package handler

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go_blog/database"
	"net/http"
	"time"
)

const endpoint = "minio.ggeta.com"
const accessKeyID = "HI4mSQabJ6GWesqES4V4"
const secreteAccessKey = "WIK6SwKqceiPCalmhDj4meOdqLdErSfw4QNpEZxx"
const bucketName = "blog-public-data"

func Upload_file(c *gin.Context, db_user *sql.DB) {
	type UploadFileRequest struct {
		FileName  string `json:"file_name"`
		HashCrc32 string `json:"hash_crc32"`
	}
	type UploadFileResponse struct {
		UploadLink string `json:"upload_link"`
		Message    string `json:"message"`
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	if user.Email == "" {
		c.JSON(http.StatusForbidden, UploadFileResponse{
			UploadLink: "",
			Message:    "permission denied",
		})
		return
	}
	// ok get request
	var uploadFileRequest UploadFileRequest
	if c.BindJSON(&uploadFileRequest) != nil {
		c.JSON(http.StatusBadRequest, UploadFileResponse{
			UploadLink: "",
			Message:    "invalid request",
		})
		return
	}
	// file_name append with hashCrc32
	file_name_with_hash := uploadFileRequest.FileName + "_" + uploadFileRequest.HashCrc32 // this should not exist in minio (skip check here)

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
		UploadLink: presignedURL.String(),
		Message:    "success",
	})
	//// upload test file with presigned url
	//// read file  /home/zong/Desktop/test.txt_12345678
	//data, err := ioutil.ReadFile("/home/zong/Desktop/test.txt_12345678")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//req, err := http.NewRequest("PUT", presignedURL.String(), bytes.NewReader(data))
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//// upload file
	//resp, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(resp)

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
