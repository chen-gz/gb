package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	//context "golang.org/x/net/context"
	"log"
	"time"
)

func connect() error {
	// add onfigure for minio
	endpoint := "minio.ggeta.com"
	accessKeyID := "DV2dmB1KJtlsP0Ud"
	secretAccessKey := "mxAiL0iQSlJiQ6lZqMBvjT261FO5mfz0"
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		fmt.Println("error")
		log.Fatalln(err)
	}
	log.Printf("%#v\n", minioClient) // minioClient is now setup
	// List all buckets
	buckets, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		fmt.Println("error")
		log.Fatalln(err)
	}
	// connect to blog backend storage
	// the bucket name is "blogBackend"

	fmt.Println(buckets)

	// check the bucket is exist or not
	// if not exist, create a new bucket
	bucketName := "blogBackend"
	exist, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return err
	}
	if !exist {
		// create a new bucket
		// set the bucket policy to public
		// set the bucket policy to read only
		err = minioClient.MakeBucket(
			context.Background(),
			bucketName,
			minio.MakeBucketOptions{Region: "us-east-1"})
	}

	// list all objects from cloudreve with a matching prefix.
	objectsCh := minioClient.ListObjects(context.Background(),
		"cloudreve",
		minio.ListObjectsOptions{Prefix: ""})

	for object := range objectsCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return object.Err
		}
		fmt.Println(object)
	}

	err = minioClient.FGetObject(context.Background(), "cloudreve", "uploads/2023/04/15/1aFoeFfT_food_fixz1.png", "/tmp/myobject", minio.GetObjectOptions{})
	// get url for object "uploads/2023/04/15/1aFoeFfT_food_fixz1.png" from bucket "cloudreve"
	// set ctx expire for 10 minutes
	ctx := context.Background()

	ctx, _ = context.WithTimeout(ctx, 10*time.Minute)
	url, err := minioClient.PresignedGetObject(context.Background(), "cloudreve", "uploads/2023/04/15/1aFoeFfT_food_fixz1.png", time.Second*60*30, nil)
	fmt.Println(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// new dynamic blog system
// markdown render which not change the connect between "$" and "$$"

func upload(filePath string) error {
	bucketName := "blogBackend"
	endpoint := "minio.ggeta.com"
	accessKeyID := "DV2dmB1KJtlsP0Ud"
	secretAccessKey := "mxAiL0iQSlJiQ6lZqMBvjT261FO5mfz0"
	useSSL := true

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println("error")
		log.Fatalln(err)
	}
	// get current time as the object name
	name := time.Now().Format("2006/01/02/15:04:05")

	minioClient.FPutObject(context.Background(),
		bucketName,
		name,
		filePath,
		minio.PutObjectOptions{})
	return nil
}

func download(filename string) error {
	bucketName := "blogBackend"
	endpoint := "minio.ggeta.com"
	// accessKeyID := "DV2dmB1KJtlsP0Ud"
	useSSL := true
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println("error")
		log.Fatalln(err)
	}

	minioClient.PresignedGetObject(
		context.Background(),
		bucketName,
		filename,
		time.Hour*10,
		nil,
	)
	return nil

}
