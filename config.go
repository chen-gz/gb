package main

import (
	"encoding/json"
	"fmt"
	"go_blog/database"
	"go_blog/handler"
	"os"
)

type Config struct {
	BlogDatabase database.BlogDbConfig `json:"blog_database"`
	UserDatabase database.UserDbConfig `json:"user_database"`
	Minio        handler.MinioConfig   `json:"minio"`
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
