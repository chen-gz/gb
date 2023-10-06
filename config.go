package main

import (
	"encoding/json"
	"fmt"
	"go_blog/database"
	"os"
)

type Minio struct {
	Endpoint         string `json:"endpoint"`
	AccessKeyID      string `json:"access_key_id"`
	SecreteAccessKey string `json:"secrete_access_key"`
	BucketName       string `json:"bucket_name"`
}
type Config struct {
	BlogDatabase database.BlogDbConfig `json:"blog_database"`
	UserDatabase database.UserDbConfig `json:"user_database"`
	Minio        Minio                 `json:"minio"`
}

// read config.json and return Config struct
func ReadConfig() Config {
	var config Config
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(configFile)
		fmt.Println(err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
