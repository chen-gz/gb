package storage
func connect() {
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
    fmt.Println(buckets)
    // list all objects from cloudreve with a matching prefix.
    objectsCh := minioClient.ListObjects(context.Background(), "cloudreve", minio.ListObjectsOptions{Prefix: ""})
    for object := range objectsCh {
        if object.Err != nil {
            fmt.Println(object.Err)
            return
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
        return
    }

}
// new dynamic blog system
// markdown render which not change the connect between "$" and "$$"
