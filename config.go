package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Database struct {
	MariadbAddress  string `json:"mariadb_address"`
	MariadbUser     string `json:"mariadb_user"`
	MariadbPassword string `json:"mariadb_password"`
}
type Minio struct {
	Endpoint         string `json:"endpoint"`
	AccessKeyID      string `json:"access_key_id"`
	SecreteAccessKey string `json:"secrete_access_key"`
	BucketName       string `json:"bucket_name"`
}
type Config struct {
	Database Database `json:"database"`
	Minio    Minio    `json:"minio"`
}

// read config.json and return Config struct
func ReadConfig() Config {
	var config Config
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
