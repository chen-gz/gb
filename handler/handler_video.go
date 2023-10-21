package handler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go_blog/database"
	"net/http"
	"time"
)

func AddVideo(c *gin.Context, dbUser *sql.DB, dbVideo *sql.DB, minioClient *minio.Client) {
	type request struct {
		Md5    string `json:"md5"`
		Sha256 string `json:"sha256"`
	}
	type response struct {
		Video        database.VideoItem `json:"video"`
		Message      string             `json:"message"`
		PresignedUrl string             `json:"presigned_url"`
	}
	var req request
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, response{
			Message: "invalid request",
		})
		return
	}
	user, ok := checkRequest(c, dbUser, &req)
	if !ok {
		return
	}
	// check md5 and sha256 length is valid or not
	if len(req.Md5) != 32 || len(req.Sha256) != 64 {
		c.JSON(http.StatusBadRequest, response{
			Message: "invalid md5 or sha256",
		})
		return
	}
	// check md5 and sha256 is exist or not
	video := database.GetVideoByMd5Sha256(dbVideo, req.Md5, req.Sha256)
	if video.Id != 0 {
		c.JSON(http.StatusBadRequest, response{
			Message: "video already exist",
		})
		return
	}
	// get presigned url
	presignedUrl, err := getPresignedUrl(minioClient, user, req.Md5, req.Sha256)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{
			Message: "get presigned url error",
		})
		return
	}
	// insert video
	video = database.VideoItem{
		UserId: user.Id,
		Md5:    req.Md5,
		Sha256: req.Sha256,
	}
	err = database.InsertVideoUser(dbVideo, video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{
			Message: "insert video error",
		})
		return
	}
	c.JSON(http.StatusOK, response{
		Video:        video,
		Message:      "success",
		PresignedUrl: presignedUrl,
	})
}
func getPresignedUrl(minioClient *minio.Client, user database.User, md5 string, sha256 string) (string, error) {
	objectName := fmt.Sprintf("%d/%s/%s", user.Id, md5[0:5], sha256[0:5])
	presignedUrl, err := minioClient.PresignedPutObject(context.Background(), VideominioConfig.BucketName, objectName, time.Hour*2)
	if err != nil {
		return "", err
	}
	return presignedUrl.String(), nil
}

var VideominioConfig MinioConfig

func InitVideoMinioClient(_config MinioConfig) *minio.Client {
	VideominioConfig = _config
	minioClient, err := minio.New(_config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(_config.AccessKeyID, _config.SecreteAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		panic(err)
	}
	return minioClient
}

func process(minioClient *minio.Client, user database.User, md5 string, sha256 string) (err error) {
	// get movie info and get m3u8
	//
	//objectName := fmt.Sprintf("%d/%s/%s", user.Id, md5[0:5], sha256[0:5])
	// download movie from minio
	//object := minio.GetObjectOptions{}

	//object, err := minioClient.GetObject(context.Background(), VideominioConfig.BucketName, objectName, minio.GetObjectOptions{})
	//if err != nil {
	//	panic(err)
	//}
	//return err

	// convert to m3u8 use ffmpeg

	// upload m3u8 to minio
	return nil
}
