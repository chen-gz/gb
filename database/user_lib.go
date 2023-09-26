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

func userTableInit() error {
	// create user table in a seperate database for user login
	db, err := sql.Open(userDbType, userDbName)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
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
		return err
	}
	// create table store the user info related to post
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS post_user_data (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
  	 	role TEXT NOT NULL,
		level INTEGER NOT NULL,
		name TEXT NOT NULL,
		group_name TEXT NOT NULL);`)
	if err != nil {
		return err
	}
	// add admin user to both tables
	_, err = db.Exec(`INSERT INTO post_user_data (email, role, level, name, group_name) VALUES (?, ?, ?, ?, ?)`,
		"admin", "admin", 100, "admin", "admin")
	return err
}
func updateUserLogin(login UserLoginData) error {
	db, err := sql.Open(userDbType, userDbName)
	if err != nil {
		return err
	}
	// update login data
	_, err = db.Exec(`UPDATE users SET email=?, password=?, user_name=? WHERE id=?`,
		login.Email, login.Password, login.UserName, login.Id)
	return err

}
func updatePUserPost(post PostUserData) error {
	db, _ := sql.Open(userDbType, userDbName)
	// update post data
	_, err := db.Exec(`UPDATE post_user_data SET email=?, role=?, level=?, name=?, group_name=? WHERE id=?`,
		post.Email, post.Role, post.Level, post.Name, post.Group, post.Id)
	return err

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

// getPostUser takes an integer ID as an input and returns the corresponding PostUserData and an error, if any.
// This function queries the post_users table in the given database to retrieve the required PostUserData entry,
// fetching the user's email, role, level, name, and group.
// In case of an error or if the user is not found, it returns an empty PostUserData struct and the error.
func getPostUser(id int) (PostUserData, error) {
	// get user data from database
	db, err := sql.Open(userDbType, userDbName)
	if err != nil {
		return PostUserData{}, err
	}
	defer db.Close()
	row := db.QueryRow(`SELECT * FROM post_users WHERE id=?`, id)
	var userdata PostUserData
	err = row.Scan(&userdata.Email, &userdata.Role, &userdata.Level, &userdata.Name, &userdata.Group)
	if err != nil {
		return PostUserData{}, err
	}
	return userdata, nil
}
