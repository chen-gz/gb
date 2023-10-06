package main

import (
	"log"
	"testing"
)

func TestReadConfig(t *testing.T) {
	// read json config file
	config := ReadConfig()
	log.Println(config.Database)
}
