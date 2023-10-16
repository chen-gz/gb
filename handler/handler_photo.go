package handler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go_blog/database"
	"log"
	"net/http"
	"net/url"
	"time"
)

func signGetPhotoLink(photo database.PhotoItemV2, client *minio.Client) (ori_url string, jpg_url string, thumb_url string, err error) {
	context := context.Background()
	var tmp *url.URL
	if photo.HasOriginal {
		filepath := fmt.Sprintf("%d_%s.%s", photo.Id, photo.OriHash[0:10], photo.OriExt)
		tmp, err = client.PresignedGetObject(context, PhotoMinioConfig.BucketName, filepath, time.Minute*10, nil)
		ori_url = tmp.String()
		if err != nil {
			log.Println("GetPhoto: ", err)
			return "", "", "", err
		}
	}
	file_path := fmt.Sprintf("%d_%s.jpg", photo.Id, photo.JpgHash[0:10])
	tmp, err = client.PresignedGetObject(context, PhotoMinioConfig.BucketName, file_path, time.Minute*10, nil)
	jpg_url = tmp.String()
	if err != nil {
		log.Println("GetPhoto: ", err)
		return "", "", "", err
	}
	file_path = fmt.Sprintf("%d_%s.jpg", photo.Id, photo.ThumbHash[0:10])
	tmp, err = client.PresignedGetObject(context, PhotoMinioConfig.BucketName, file_path, time.Minute*10, nil)
	thumb_url = tmp.String()
	if err != nil {
		log.Println("GetPhoto: ", err)
		return "", "", "", err
	}
	return ori_url, jpg_url, thumb_url, nil
}

func signPutPhotoLink(photo database.PhotoItemV2, client *minio.Client) (ori_url string, jpg_url string, thumb_url string, err error) {
	ctx := context.Background()
	if photo.HasOriginal {
		filepath := fmt.Sprintf("%d_%s.%s", photo.Id, photo.OriHash[0:10], photo.OriExt)
		tmp, err := client.PresignedPutObject(ctx, PhotoMinioConfig.BucketName, filepath, time.Minute*10)
		ori_url = tmp.String()
		if err != nil {
			log.Println("GetPhoto: ", err)
			return "", "", "", err
		}
	}
	file_path := fmt.Sprintf("%d_%s.jpg", photo.Id, photo.JpgHash[0:10])
	tmp, err := client.PresignedPutObject(ctx, PhotoMinioConfig.BucketName, file_path, time.Minute*10)
	jpg_url = tmp.String()
	if err != nil {
		log.Println("GetPhoto: ", err)
		return "", "", "", err
	}
	file_path = fmt.Sprintf("%d_%s.jpg", photo.Id, photo.ThumbHash[0:10])
	tmp, err = client.PresignedPutObject(ctx, PhotoMinioConfig.BucketName, file_path, time.Minute*10)
	thumb_url = tmp.String()
	if err != nil {
		log.Println("GetPhoto: ", err)
		return "", "", "", err
	}
	return ori_url, jpg_url, thumb_url, nil
}

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
func GetDeletedPhotoIds(c *gin.Context, db_user *sql.DB, db_photo *sql.DB) {
	type GetDeletedPhotoIdsResponse struct {
		Ids     []int  `json:"ids"`
		Message string `json:"message"`
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	if user.Id == 0 {
		c.JSON(http.StatusForbidden, GetDeletedPhotoIdsResponse{Message: "permission denied"})
		return
	}
	var getDeletedPhotoIdsResponse GetDeletedPhotoIdsResponse
	ids, err := database.GetDeletedPhotoList(db_photo, user)
	if err != nil {
		log.Println("GetDeletedPhotoIds Error: ", err)
		c.JSON(http.StatusNotFound, GetDeletedPhotoIdsResponse{Message: "not found"})
		return
	}
	getDeletedPhotoIdsResponse.Ids = ids
	getDeletedPhotoIdsResponse.Message = "ok"
	c.JSON(http.StatusOK, getDeletedPhotoIdsResponse)
}

func UpdatePhoto(c *gin.Context, db_user *sql.DB, db_photo *sql.DB) {
	type UpdatePhotoRequest database.PhotoItem
	type UpdatePhotoResponse struct {
		Message string `json:"message"`
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	if user.Id == 0 {
		c.JSON(http.StatusForbidden, UpdatePhotoResponse{Message: "permission denied"})
		return
	}
	var updatePhotoRequest UpdatePhotoRequest
	if c.BindJSON(&updatePhotoRequest) != nil {
		c.JSON(http.StatusBadRequest, UpdatePhotoResponse{Message: "invalid request"})
		return
	}
	err := database.UpdatePhotoUser(db_photo, user, database.PhotoItem(updatePhotoRequest))
	if err != nil {
		log.Println("UpdatePhoto Error: ", err)
		c.JSON(http.StatusNotFound, UpdatePhotoResponse{Message: "not found"})
		return
	}
	c.JSON(http.StatusOK, UpdatePhotoResponse{Message: "ok"})
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

func GetPhotoHash(c *gin.Context, hash string, db_user *sql.DB, db_photo *sql.DB, client *minio.Client) {
	type GetPhotoResponse struct {
		Photo   database.PhotoItemV2 `json:"photo"`
		ThumUrl string               `json:"thum_url"`
		OriUrl  string               `json:"ori_url"`
		JpegUrl string               `json:"jpeg_url"`
		Message string               `json:"message"`
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	if user.Id == 0 {
		c.JSON(http.StatusForbidden, GetPhotoResponse{Message: "permission denied"})
		return
	}

	var getPhotoResponse GetPhotoResponse
	photo, err := database.GetPhotoByJpgHash(db_photo, user, hash)
	if err != nil {
		log.Println("GetPhoto Error: ", err)
		c.JSON(http.StatusNotFound, GetPhotoResponse{Message: "hash not found"})
		return
	}
	getPhotoResponse.Photo = photo
	getPhotoResponse.Message = "ok"
	getPhotoResponse.OriUrl, getPhotoResponse.JpegUrl, getPhotoResponse.ThumUrl, err = signGetPhotoLink(photo, client)
	if err != nil {
		log.Println("GetPhoto: ", err)
		c.JSON(http.StatusInternalServerError, GetPhotoResponse{Message: "internal error"})
		return
	}
	c.JSON(http.StatusOK, getPhotoResponse)
}

func GetPhotoId(c *gin.Context, id int, db_user *sql.DB, db_photo *sql.DB, client *minio.Client) {
	type GetPhotoResponse struct {
		Photo   database.PhotoItemV2 `json:"photo"`
		ThumUrl string               `json:"thum_url"`
		OriUrl  string               `json:"ori_url"`
		JpegUrl string               `json:"jpeg_url"`
		Message string               `json:"message"`
	}
	user := database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	if user.Id == 0 {
		c.JSON(http.StatusForbidden, GetPhotoResponse{Message: "permission denied"})
		return
	}
	photo, err := database.GetPhotoById(db_photo, user, id)
	if err != nil {
		log.Println("GetPhoto Error: ", err)
		c.JSON(http.StatusNotFound, GetPhotoResponse{Message: "id not found"})
		return
	}
	var getPhotoResponse GetPhotoResponse
	getPhotoResponse.Photo = photo
	getPhotoResponse.Message = "ok"
	getPhotoResponse.OriUrl, getPhotoResponse.JpegUrl, getPhotoResponse.ThumUrl, err = signGetPhotoLink(photo, client)
	if err != nil {
		log.Println("GetPhoto: ", err)
		c.JSON(http.StatusInternalServerError, GetPhotoResponse{Message: "internal error"})
		return
	}
	c.JSON(http.StatusOK, getPhotoResponse)

}

func UpdatePhotoMeta(c *gin.Context, db_user *sql.DB, db_photo *sql.DB) {
	// only some fields can be updated
	// tag, category, deleted
	type UpdatePhotoMetaRequest struct {
		Id       int    `json:"id"`
		Deleted  bool   `json:"deleted"`
		Tag      string `json:"tag"`
		Category string `json:"category"`
	}
	type UpdatePhotoMetaResponse struct {
		Message string `json:"message"`
	}
	var updatePhotoMetaRequest UpdatePhotoMetaRequest
	user, ok := checkRequest(c, db_user, &updatePhotoMetaRequest)
	if !ok {
		return
	}

	photo, err := database.GetPhotoById(db_photo, user, updatePhotoMetaRequest.Id)
	if err != nil {
		log.Println("UpdatePhotoMeta: ", err)
		c.JSON(http.StatusNotFound, UpdatePhotoMetaResponse{Message: "id not found"})
		return
	}
	photo.Tags = updatePhotoMetaRequest.Tag
	photo.Category = updatePhotoMetaRequest.Category
	err = database.UpdatePhotoMetaById(db_photo, user, photo)
	if err != nil {
		log.Println("UpdatePhotoMeta: ", err)
		c.JSON(http.StatusInternalServerError, UpdatePhotoMetaResponse{Message: "internal error"})
		return
	}
	c.JSON(http.StatusOK, UpdatePhotoMetaResponse{Message: "ok"})
}

func UpdatePhotoFile(c *gin.Context, dbUser *sql.DB, dbPhoto *sql.DB, client *minio.Client) {
	type UpdatePhotoFileRequest database.PhotoItemV2
	type UpdatePhotoFileResponse struct {
		Id                   int    `json:"id"`
		Message              string `json:"message"`
		PresignedOriginalUrl string `json:"presigned_original_url"`
		PresignedJpegUrl     string `json:"presigned_jpeg_url"`
		PresignedThumbUrl    string `json:"presigned_thumb_url"`
	}
	var updatePhotoFileRequest UpdatePhotoFileRequest
	var updatePhotoFileResponse UpdatePhotoFileResponse
	user, ok := checkRequest(c, dbUser, &updatePhotoFileRequest)
	if !ok {
		return
	}
	photo, err := database.GetPhotoById(dbPhoto, user, updatePhotoFileRequest.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, UpdatePhotoFileResponse{Message: "id not found"})
		return
	}
	log.Println("UpdatePhotoFile: ", updatePhotoFileRequest)
	if photo.OriHash != updatePhotoFileRequest.OriHash && len(updatePhotoFileRequest.OriHash) > 0 {
		// delete old original file and upload new original file
		fmt.Println("%+v", photo)
		if photo.OriHash != "" && photo.OriExt == "" {

			oldFilePath := fmt.Sprintf("%d_%s.%s", photo.Id, photo.OriHash[0:10], photo.OriExt)
			err = client.RemoveObject(context.Background(), PhotoMinioConfig.BucketName, oldFilePath, minio.RemoveObjectOptions{})
			if err != nil {
				log.Println("UpdatePhotoFile remove old file failed: ", err)
				//c.JSON(http.StatusInternalServerError, UpdatePhotoFileResponse{Message: "internal error"})
				//return
			}
		}
		newFilePath := fmt.Sprintf("%d_%s.%s", photo.Id, updatePhotoFileRequest.OriHash[0:10], updatePhotoFileRequest.OriExt)
		photo.OriHash = updatePhotoFileRequest.OriHash
		photo.OriExt = updatePhotoFileRequest.OriExt
		photo.HasOriginal = true
		url, err := client.PresignedPutObject(c, PhotoMinioConfig.BucketName, newFilePath, time.Minute*10)
		if err != nil {
			log.Println("UpdatePhotoFile: ", err)
			c.JSON(http.StatusBadRequest, UpdatePhotoFileResponse{Message: "internal request"})
			return
		}
		updatePhotoFileResponse.PresignedOriginalUrl = url.String()
	}
	if photo.JpgHash != updatePhotoFileRequest.JpgHash && len(updatePhotoFileRequest.JpgHash) > 0 {
		// delete old jpg file and upload new jpg file
		oldFilePath := fmt.Sprintf("%d_%s.jpg", photo.Id, photo.JpgHash[0:10])
		err = client.RemoveObject(context.Background(), PhotoMinioConfig.BucketName, oldFilePath, minio.RemoveObjectOptions{})
		if err != nil {
			log.Println("UpdatePhotoFile remove old file failed: ", err)
			//log.Println("UpdatePhotoFile: ", err)
			//c.JSON(http.StatusInternalServerError, UpdatePhotoFileResponse{Message: "internal error"})
			//return
		}
		photo.JpgHash = updatePhotoFileRequest.JpgHash
		newFilePath := fmt.Sprintf("%d_%s.jpg", photo.Id, updatePhotoFileRequest.JpgHash[0:10])
		url, err := client.PresignedPutObject(c, PhotoMinioConfig.BucketName, newFilePath, time.Minute*10)
		if err != nil {
			log.Println("UpdatePhotoFile: ", err)
			c.JSON(http.StatusBadRequest, UpdatePhotoFileResponse{Message: "invalid request"})
			return
		}
		updatePhotoFileResponse.PresignedJpegUrl = url.String()
	}
	if photo.ThumbHash != updatePhotoFileRequest.ThumbHash && len(updatePhotoFileRequest.ThumbHash) > 0 {
		// delete old thumbnail file and upload new thumbnail file
		oldFilePath := fmt.Sprintf("%d_%s.jpg", photo.Id, photo.ThumbHash[0:10])
		err = client.RemoveObject(context.Background(), PhotoMinioConfig.BucketName, oldFilePath, minio.RemoveObjectOptions{})
		//photo.OriHash+"_thumbnail.jpg", minio.RemoveObjectOptions{})
		if err != nil {
			//log.Println("UpdatePhotoFile: ", err)
			//c.JSON(http.StatusInternalServerError, UpdatePhotoFileResponse{Message: "internal error"})
			//return
			log.Println("UpdatePhotoFile remove old file failed: ", err)
		}
		photo.ThumbHash = updatePhotoFileRequest.ThumbHash
		newFilePath := fmt.Sprintf("%d_%s.jpg", photo.Id, updatePhotoFileRequest.ThumbHash[0:10])
		url, err := client.PresignedPutObject(c, PhotoMinioConfig.BucketName, newFilePath, time.Minute*10)
		//updatePhotoFileRequest.OriHash+"_thumbnail.jpg", time.Hour)
		if err != nil {
			log.Println("UpdatePhotoFile: ", err)
			c.JSON(http.StatusBadRequest, UpdatePhotoFileResponse{Message: "invalid request"})
			return
		}
		updatePhotoFileResponse.PresignedThumbUrl = url.String()
	}
	// update database

	err = database.UpdatePhotoById(dbPhoto, user, photo)
	if err != nil {
		log.Println("UpdatePhotoFile: ", err)
		c.JSON(http.StatusInternalServerError, UpdatePhotoFileResponse{Message: "internal error"})
		return
	}
	updatePhotoFileResponse.Id = photo.Id
	updatePhotoFileResponse.Message = "ok"
	c.JSON(http.StatusOK, updatePhotoFileResponse)

}
func InsertPhotoV2(c *gin.Context, dbUser *sql.DB, dbPhoto *sql.DB, client *minio.Client) {
	type InsertPhotoRequest database.PhotoItemV2 // id will be ignored
	type InsertPhotoResponse struct {
		Id                   int    `json:"id"`
		Message              string `json:"message"`
		PresignedOriginalUrl string `json:"presigned_original_url"`
		PresignedJpegUrl     string `json:"presigned_jpeg_url"`
		PresignedThumbUrl    string `json:"presigned_thumb_url"`
	}
	var insertPhotoResponse InsertPhotoResponse
	var insertPhotoRequest InsertPhotoRequest
	user, ok := checkRequest(c, dbUser, &insertPhotoRequest)
	if !ok {
		return
	}
	photo, err := database.InsertPhotoUserV2(dbPhoto, user, database.PhotoItemV2(insertPhotoRequest))
	// presign url asd return
	insertPhotoResponse.Message = "ok"

	if err != nil {
		log.Println("InsertPhoto: ", err)
		c.JSON(http.StatusBadRequest, InsertPhotoResponse{Message: fmt.Sprintf("%s", err)})
		return
	}
	insertPhotoResponse.PresignedOriginalUrl, insertPhotoResponse.PresignedJpegUrl, insertPhotoResponse.PresignedThumbUrl, err = signPutPhotoLink(photo, client)
	insertPhotoResponse.Id = photo.Id
	if err != nil {
		log.Println("InsertPhoto: ", err)
		c.JSON(http.StatusInternalServerError, InsertPhotoResponse{Message: "internal error"})
		return
	}
	c.JSON(http.StatusOK, insertPhotoResponse)
}

// func (c *Context) BindJSON(obj any) error {
func checkRequest(c *gin.Context, db_user *sql.DB, request any) (user database.User, ok bool) {
	user = database.V3GetUserByAuthHeader(db_user, c.Request.Header.Get("Authorization"))
	if user.Id == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return user, false
	}
	if c.BindJSON(&request) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return user, false
	}
	return user, true
}
