package database

import "testing"

func TestCreateDatabase(t *testing.T) {
	// create database
	err := userTableInit()
	if err != nil {
		t.Error(err)
	}
}
