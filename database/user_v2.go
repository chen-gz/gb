package database

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

// design new user system
type Users struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Group string `json:"group"`
	Pass  string `json:"pass"`
}

func Login(email string, pass string) {
	// query the database
	db, err := sql.Open(userDbType, userDbName)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(`SELECT id FROM users WHERE email=? AND password=?`, email, pass)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return
		}
		return
	}
	return
}

func AddUser(users Users) {
	db, err := sql.Open(userDbType, userDbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(`INSERT INTO users (email, password, user_name) VALUES (?, ?, ?)`,
		users.Email, users.Pass, users.Name)
	if err != nil {
		return
	}
	// email should be unique
}

func Deleteuser(email string) {
	db, err := sql.Open(userDbType, userDbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(`DELETE FROM users WHERE email=?`, email)
	if err != nil {
		return
	}
}

func UpdateUser(users Users) {

}

var secreteKey = []byte("bcb967bec859b86e96564992792636bb442548af35a2e3374cee7a0f92542c18")

func V1GenerateToken(email string) string {
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

func GetUserByAuthHeader(info string) db.UserData {
	// if auth type is Bearer get token
	if len(info) < 7 {
		return db.UserData{}
	}
	if info[0:7] == "Bearer " {
		token := info[7:]
		valid, email := V1VerifyToken(token)
		if !valid {
			return db.UserData{}
		} else {
			return db.GetUser(email)
		}
	} else if info[0:6] == "Basic " {
		// todo: handle invalid auth header.
		userPass := info[6:]
		log.Println(userPass)
		email := userPass[0:strings.Index(userPass, ":")]
		pass := userPass[strings.Index(userPass, ":")+1:]
		// get
		if login(email, pass) {
			return db.GetUser(email)
		}
		return db.UserData{}
	}
	return db.UserData{}
}

func login(email string, password string) bool {
	if email == "chen-gz@outlook.com" && password == "Connie" {
		return true
	}
	return false
}

func V2Login(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	user := GetUserByAuthHeader(auth)
	if user.Email != "" && auth[0:6] == "Basic " {
		c.JSON(http.StatusOK, gin.H{
			"email": user.Email,
			"token": V1GenerateToken(user.Email),
			"msg":   "log in success",
		})
	} else if user.Email != "" && auth[0:7] == "Bearer " {
		c.JSON(http.StatusOK, gin.H{
			"email": user.Email,
			"token": auth[7:],
			"msg":   "log in success",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"email": "",
			"token": "",
			"msg":   "log in failed",
		})

	}
}
