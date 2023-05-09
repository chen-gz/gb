package database

//package main

type UserData struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	Level int    `json:"level"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

func GetUser(email string) UserData {
	// get user data from database
	if email == "chen-gz@outlook.com" {
		return UserData{
			Email: "chen-gz@outlook.com",
			Role:  "admin",
			Level: 0,
			Name:  "Guangzong",
			Group: "admin",
		}
	}
	return UserData{}
}
