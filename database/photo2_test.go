package database

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type MinioConfig struct {
	Endpoint         string `json:"endpoint"`
	AccessKeyID      string `json:"access_key_id"`
	SecreteAccessKey string `json:"secrete_access_key"`
	BucketName       string `json:"bucket_name"`
}

var config MinioConfig

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

// migration from v1 to v2
//func TestMigration(t *testing.T) {
//	// get photo from v1
//
//	config := PhotoDbConfig{
//		Address:       "tcp(192.168.0.174:3306)",
//		User:          "zong",
//		Password:      "Connie",
//		PhotoDatabase: "eta_photo",
//	}
//	minio_config := MinioConfig{
//		Endpoint:         "minio.ggeta.com",
//		AccessKeyID:      "HI4mSQabJ6GWesqES4V4",
//		SecreteAccessKey: "WIK6SwKqceiPCalmhDj4meOdqLdErSfw4QNpEZxx",
//		BucketName:       "photo",
//	}
//
//	db_photo, _ := InitPhotoDb(config)
//	minio_client := InitMinioClient(minio_config)
//	user := User{Id: 2, Name: "Guangzong Chen"}
//	InitPhotoTableV2(db_photo, user)
//	ids, _ := GetAllPhotoList(db_photo, user)
//	//log.Println(ids, err)
//	// get photo
//	for _, id := range ids {
//		var oriFile, jpgFile, thumbFile minio.ObjectInfo
//		//startTime := time.Now()
//		photo, err := GetPhotoUser(db_photo, user, id)
//		//fmt.Println("Get photo: ", time.Since(startTime))
//		//startTime = time.Now()
//
//		//log.Println(photo, err)
//		if photo.HasOriginal {
//			// check files
//			fileName := fmt.Sprintf("%s_ori.%s", photo.Hash, photo.OriginalExt)
//			//log.Println(fileName)
//			oriFile, err = minio_client.StatObject(context.Background(), minio_config.BucketName, fileName, minio.StatObjectOptions{})
//			if err != nil {
//				//log.Println(err)
//				log.Println(fileName, err)
//			}
//		}
//		//fmt.Println("Get ori file: ", time.Since(startTime))
//		//startTime = time.Now()
//		fileName := fmt.Sprintf("%s.jpg", photo.Hash)
//		jpgFile, err = minio_client.StatObject(context.Background(), minio_config.BucketName, fileName, minio.StatObjectOptions{})
//		if err != nil {
//			//log.Println(err)
//			log.Println(fileName, err)
//		}
//		//fmt.Println("Get jpg file: ", time.Since(startTime))
//		//startTime = time.Now()
//		fileName = fmt.Sprintf("%s_thumbnail.jpg", photo.Hash)
//		thumbFile, err = minio_client.StatObject(context.Background(), minio_config.BucketName, fileName, minio.StatObjectOptions{})
//		ok := true
//		if err != nil {
//			//log.Println(err)
//			log.Println(fileName, err)
//			ok = false
//		}
//		//fmt.Println("Get thumb file: ", time.Since(startTime))
//		//startTime = time.Now()
//		_, _, _ = oriFile, jpgFile, thumbFile
//		photoV2 := PhotoItemV2{
//			Id:          photo.Id,
//			OriHash:     oriFile.ETag,
//			JpgHash:     jpgFile.ETag,
//			ThumbHash:   thumbFile.ETag,
//			HasOriginal: photo.HasOriginal,
//			OriExt:      photo.OriginalExt,
//			Deleted:     photo.Deleted,
//			Tags:        photo.Tags,
//			Category:    photo.Category,
//		}
//
//		photoV22, err := InsertPhotoUserV2(db_photo, user, photoV2)
//		//log.Println(photoV22, err)
//		if err != nil {
//			log.Println(photoV22, err)
//			//log.Println(err)
//		}
//		//fmt.Println("Insert photo: ", time.Since(startTime))
//		//copy files
//		//minio_client.RenameObject(context.Background(), minio_config.BucketName, fmt.Sprintf("%s.jpg", photo.Hash), fmt.Sprintf("%s.jpg", photoV22.JpgHash), minio.CopySrcOptions{})
//
//		if photo.HasOriginal {
//			oldFileName := fmt.Sprintf("%s_ori.%s", photo.Hash, photo.OriginalExt)
//			newFileName := fmt.Sprintf("%d_%s.%s", photoV22.Id, photoV22.OriHash[0:10], photoV22.OriExt)
//			src := minio.CopySrcOptions{
//				Bucket: minio_config.BucketName,
//				Object: oldFileName,
//			}
//			dst := minio.CopyDestOptions{
//				Bucket: minio_config.BucketName,
//				Object: newFileName,
//			}
//			_, err2 := minio_client.CopyObject(context.Background(), dst, src)
//			if err2 != nil {
//				log.Println(err2)
//			}
//		}
//		// copy jpg
//		oldFileName := fmt.Sprintf("%s.jpg", photo.Hash)
//		newFileName := fmt.Sprintf("%d_%s.jpg", photoV22.Id, photoV22.JpgHash[0:10])
//		src := minio.CopySrcOptions{
//			Bucket: minio_config.BucketName,
//			Object: oldFileName,
//		}
//		dst := minio.CopyDestOptions{
//			Bucket: minio_config.BucketName,
//			Object: newFileName,
//		}
//		_, err = minio_client.CopyObject(context.Background(), dst, src)
//		if err != nil {
//			//return
//			log.Println(err)
//		}
//		// copy thumb
//		if ok {
//			oldFileName = fmt.Sprintf("%s_thumbnail.jpg", photo.Hash)
//			newFileName = fmt.Sprintf("%d_%s.jpg", photoV22.Id, photoV22.ThumbHash[0:10])
//			src = minio.CopySrcOptions{
//				Bucket: minio_config.BucketName,
//				Object: oldFileName,
//			}
//			dst = minio.CopyDestOptions{
//				Bucket: minio_config.BucketName,
//				Object: newFileName,
//			}
//			_, err2 := minio_client.CopyObject(context.Background(), dst, src)
//			if err2 != nil {
//				//return
//				log.Println(err2)
//			}
//		}
//		//fmt.Println("Copy files: ", src, dst)
//		//break
//	}
//
//}
