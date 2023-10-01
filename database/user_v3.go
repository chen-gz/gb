package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

func UserDbInit() (db_user *sql.DB, err error) {
	db, err := sql.Open("mysql", "zong:Connie@tcp(192.168.0.174:3306)/")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS eta_user")

	db_user, err = sql.Open("mysql", "zong:Connie@tcp(192.168.0.174:3306)/eta_user")
	if err != nil {
		panic(err)
	}
	_, err = db_user.Exec(
		` CREATE TABLE IF NOT EXISTS users (
    id         INT UNSIGNED AUTO_INCREMENT,
    email      VARCHAR(255) UNIQUE NOT NULL,
    name       VARCHAR(255),
    password   VARCHAR(255),
    created_at TIMESTAMP                                                      DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP                                                      DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);`)

	if err != nil {
		panic(err)
	}
	return db_user, nil
}

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func UserAdd(db_user *sql.DB, user User, password string) error {
	_, err := db_user.Exec("INSERT INTO users (email, name, password) VALUES (?, ?, ?)", user.Email, user.Name, password)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

// GetUserByEmail get user by email
// if cannot find user, return empty user
func GetUserByEmail(dbUser *sql.DB, email string) User {
	var user User
	err := dbUser.QueryRow("SELECT id, email, name FROM users WHERE email=?", email).Scan(&user.Id, &user.Email, &user.Name)
	if err != nil {
		log.Fatalln(err)
		return User{}
	}
	return user
}

var secreteKey = []byte("bcb967bec859b86e96564992792636bb442548af35a2e3374cee7a0f92542c18")

func V1VerifyToken(token string) (bool, string) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secreteKey, nil
	})
	if err != nil {
		log.Println("verify token failed: ", err)
		return false, ""
	}
	valid := parsedToken.Valid
	email := parsedToken.Claims.(jwt.MapClaims)["email"].(string)
	return valid, email
}

// V3GetUserByAuthHeader get user by auth header
// if auth type is Bearer, get token and verify it
// if auth type is Basic, return empty user
// if auth header is invalid, return empty user
func V3GetUserByAuthHeader(db_user *sql.DB, auth string) User {
	// if auth type is Bearer get token
	if len(auth) < 7 {
		return User{}
	}
	if auth[0:7] == "Bearer " {
		token := auth[7:]
		valid, email := V1VerifyToken(token)
		if !valid {
			return User{}
		} else {
			return GetUserByEmail(db_user, email)
		}
	}
	return User{}
}

func V3Login(db_user *sql.DB, email string, password string) bool {
	// select rwo from users where email = email and password = password
	err := db_user.QueryRow("SELECT email FROM users WHERE email=? AND password=?", email, password).Scan(&email)
	log.Println("login result: ", err)
	if err != nil {
		return false
	}
	return true
}

func V3GenerateToken(email string) string {
	log.Println("Generating token for user: ", email, " ...")
	signingMethod := jwt.SigningMethodHS256 // HS256 is an instance of HMAC
	claims := jwt.MapClaims{
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	token := jwt.NewWithClaims(signingMethod, claims)
	signedToken, err := token.SignedString(secreteKey)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Generated token success")
	return signedToken
}
