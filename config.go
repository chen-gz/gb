package main

import (
	"encoding/json"
	"fmt"
	"go_blog/database"
	"go_blog/handler"
	"go_blog/interfaces"
	"os"
)

type Config struct {
	BlogDatabase  database.BlogDbConfig  `json:"blog_database"`
	UserDatabase  database.UserDbConfig  `json:"user_database"`
	PhotoDatabase database.PhotoDbConfig `json:"photo_database"`
	Minio         handler.MinioConfig    `json:"minio"`
	PhotoMinio    handler.MinioConfig    `json:"photo_minio"`
	VideoMinio    handler.MinioConfig    `json:"video_minio"`
	VideoDb       interfaces.DbConfig    `json:"video_db"`
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
