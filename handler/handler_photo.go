package handler

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go_blog/database"
	"log"
	"net/http"
	"time"
)

func InsertPhoto(c *gin.Context, db_user *sql.DB, db_photo *sql.DB, client *minio.Client) {
	type InsertPhotoRequest struct {
		Hash        string `json:"hash"` //  sha1 hash. If original file exist. Use the hash of the original file. If not, use the hash of the Jpeg file.
		HasOriginal bool   `json:"has_original"`
		OriginalExt string `json:"original_ext"`
	}
	type InsertPhotoResponse struct {
		Message               string `json:"message"`
		PresignedOriginalUrl  string `json:"presigned_original_url"`
		PresignedThumbnailUrl string `json:"presigned_thumbnail_url"`
		PresignedJpegUrl      string `json:"presigned_jpeg_url"`
	}
	var insertPhotoRequest InsertPhotoRequest
	var insertPhotoResponse InsertPhotoResponse
	if c.BindJSON(&insertPhotoRequest) != nil {
		c.JSON(http.StatusBadRequest, InsertPhotoResponse{Message: "invalid request - bind json error"})
		return
	}
	log.Println("InsertPhoto: ", insertPhotoRequest)
	if insertPhotoRequest.HasOriginal && insertPhotoRequest.OriginalExt == "" {
		c.JSON(http.StatusBadRequest, InsertPhotoResponse{Message: "invalid request - original ext is empty"})
		return
	}
	if insertPhotoRequest.Hash == "" {
		c.JSON(http.StatusBadRequest, InsertPhotoResponse{Message: "invalid request - hash is empty"})
		return
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	if user.Id == 0 {
		c.JSON(http.StatusForbidden, InsertPhotoResponse{Message: "permission denied"})
		return
	}
	photo := database.PhotoItem{
		Hash:        insertPhotoRequest.Hash,
		HasOriginal: insertPhotoRequest.HasOriginal,
		OriginalExt: insertPhotoRequest.OriginalExt,
		Tags:        "",
		Category:    "",
	}
	err := database.InsertPhotoUser(db_photo, user, photo)
	if err != nil {
		c.JSON(http.StatusBadRequest, InsertPhotoResponse{Message: fmt.Sprintf("InsertPhoto: ", err)})
		return
	}
	// presign url adn return
	insertPhotoResponse.Message = "ok"
	if insertPhotoRequest.HasOriginal {
		url, err := client.PresignedPutObject(c, PhotoMinioConfig.BucketName, insertPhotoRequest.Hash+"_ori."+insertPhotoRequest.OriginalExt, time.Hour)
		if err != nil {
			log.Println("InsertPhoto: ", err)
			c.JSON(http.StatusBadRequest, InsertPhotoResponse{Message: "invalid request"})
			return
		}
		insertPhotoResponse.PresignedOriginalUrl = url.String()
	}
	url, err := client.PresignedPutObject(c, PhotoMinioConfig.BucketName, insertPhotoRequest.Hash+"_thumbnail.jpg", time.Hour)
	if err != nil {
		log.Println("InsertPhoto: ", err)
		c.JSON(http.StatusBadRequest, InsertPhotoResponse{Message: "invalid request"})
		return
	}
	insertPhotoResponse.PresignedThumbnailUrl = url.String()

	url, err = client.PresignedPutObject(c, PhotoMinioConfig.BucketName, insertPhotoRequest.Hash+".jpg", time.Hour)
	if err != nil {
		log.Println("InsertPhoto: ", err)
		c.JSON(http.StatusBadRequest, InsertPhotoResponse{Message: "invalid request"})
		return
	}
	insertPhotoResponse.PresignedJpegUrl = url.String()

	c.JSON(http.StatusOK, insertPhotoResponse)
}

func GetPhoto(c *gin.Context, db_user *sql.DB, db_photo *sql.DB, client *minio.Client) {
	type GetPhotoRequest struct {
		Id int `json:"id"`
	}
	type GetPhotoResponse struct {
		Photo   database.PhotoItem `json:"photo"`
		ThumUrl string             `json:"thum_url"`
		OriUrl  string             `json:"ori_url"`
		JpegUrl string             `json:"jpeg_url"`
		Message string             `json:"message"`
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	if user.Id == 0 {
		c.JSON(http.StatusForbidden, GetPhotoResponse{Message: "permission denied"})
		return
	}
	var getPhotoRequest GetPhotoRequest
	if c.BindJSON(&getPhotoRequest) != nil {
		c.JSON(http.StatusBadRequest, GetPhotoResponse{Message: "invalid request"})
		return
	}

	var getPhotoResponse GetPhotoResponse
	photo, err := database.GetPhotoUser(db_photo, user, getPhotoRequest.Id)
	if err != nil {
		log.Println("GetPhoto Error: ", err)
		c.JSON(http.StatusNotFound, GetPhotoResponse{Message: "not found"})
		return
	}
	getPhotoResponse.Photo = photo
	getPhotoResponse.Message = "ok"
	if photo.HasOriginal {
		url, err := client.PresignedGetObject(c, PhotoMinioConfig.BucketName, photo.Hash+"_ori."+photo.OriginalExt, time.Minute*10, nil)
		if err != nil {
			log.Println("GetPhoto: ", err)
			c.JSON(http.StatusBadRequest, GetPhotoResponse{Message: "invalid request"})
			return
		}
		getPhotoResponse.OriUrl = url.String()
	}
	url, err := client.PresignedGetObject(c, PhotoMinioConfig.BucketName, photo.Hash+"_thumbnail.jpg", time.Minute*10, nil)
	if err != nil {
		log.Println("GetPhoto: ", err)
		c.JSON(http.StatusBadRequest, GetPhotoResponse{Message: "invalid request"})
		return
	}
	getPhotoResponse.ThumUrl = url.String()
	url, err = client.PresignedGetObject(c, PhotoMinioConfig.BucketName, photo.Hash+".jpg", time.Minute*10, nil)
	if err != nil {
		log.Println("GetPhoto: ", err)
		c.JSON(http.StatusBadRequest, GetPhotoResponse{Message: "invalid request"})
		return
	}
	getPhotoResponse.JpegUrl = url.String()
	log.Println("GetPhoto: ", getPhotoResponse)
	c.JSON(http.StatusOK, getPhotoResponse)
}
func GetPhotoIds(c *gin.Context, db_user *sql.DB, db_photo *sql.DB) {
	type GetPhotoIdsResponse struct {
		Ids     []int  `json:"ids"`
		Message string `json:"message"`
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	if user.Id == 0 {
		c.JSON(http.StatusForbidden, GetPhotoIdsResponse{Message: "permission denied"})
		return
	}
	var getPhotoIdsResponse GetPhotoIdsResponse
	ids, err := database.GetAllPhotoList(db_photo, user)
	if err != nil {
		log.Println("GetPhotoIds Error: ", err)
		c.JSON(http.StatusNotFound, GetPhotoIdsResponse{Message: "not found"})
		return
	}
	getPhotoIdsResponse.Ids = ids
	getPhotoIdsResponse.Message = "ok"
	c.JSON(http.StatusOK, getPhotoIdsResponse)
}

var PhotoMinioConfig MinioConfig

func InitPhotoMinioClient(_config MinioConfig) *minio.Client {
	PhotoMinioConfig = _config
	client, err := minio.New(PhotoMinioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(PhotoMinioConfig.AccessKeyID, PhotoMinioConfig.SecreteAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatal("InitMinioClient: ", err)
		return nil
	}
	return client
}
