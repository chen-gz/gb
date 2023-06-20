package database

import (
	"database/sql"
)

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
			Level: 100,
			Name:  "Guangzong",
			Group: "admin",
		}
	}
	return UserData{}
}

type UserLoginData struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"user_name"`
}

var userDbType = "sqlite3"
var userDbName = "user.db"

func user_table_init() {
	// create user table in a seperate database
	db, err := sql.Open(userDbType, userDbName)
	db.Exec(`CREATE TABLE IF NOT EXISTS users (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		email TEXT UNIQUE NOT NULL,
    		password TEXT UNIQUE  NOT NULL,
    		user_name TEXT UNIQUE NOT NULL);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT INTO users (email, password, user_name) VALUES (?, ?, ?)`,
		"admin", "admin", "admin")
	if err != nil {
		return
	}
}
func addUesr(email string, password string, userName string) {
	db, err := sql.Open(userDbType, userDbName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO users (email, password, user_name) VALUES (?, ?, ?)`,
		email, password, userName)
	if err != nil {
		return
	}
}

// return the user id if the user is valid and verified
func verifyUser(email string, password string) int {
	db, err := sql.Open(userDbType, userDbName)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(`SELECT id FROM users WHERE email=? AND password=?`, email, password)
	if err != nil {
		return -1
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return -1
		}
		return id
	}
	return -1
}

type PostUserData struct {
	Email string `json:"email"`
	Role  string `json:"role"` // author, editor, admin, guest
	Level int    `json:"level"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

func getPostUser(email string,
	password string,
	verified bool,
) (PostUserData, error) {
	if !verified && verifyUser(email, password) == -1 {
		return PostUserData{}, nil
	}
	// get user data from database
	db, err := sql.Open(userDbType, userDbName)
	defer db.Close()
	row := db.QueryRow(`SELECT * FROM post_users WHERE email=?`, email)
	var userdata PostUserData
	err = row.Scan(&userdata.Email, &userdata.Role, &userdata.Level, &userdata.Name, &userdata.Group)
	if err != nil {
		return PostUserData{}, err
	}
	return userdata, nil
}
