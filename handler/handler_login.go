//package main

package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	db "go_blog/database"
	"log"
	"net/http"
	"strings"
	"time"
)

var secreteKey = []byte("bcb967bec859b86e96564992792636bb442548af35a2e3374cee7a0f92542c18")

func V1Login(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")

	if email == "chen-gz@outlook.com" && password == "Connie" {
		c.JSON(http.StatusOK, gin.H{
			"token": V1GenerateToken(email),
			"msg":   "log in success",
			"name":  "test",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "log in failed",
		})
	}
}
func V1Verify(c *gin.Context) {
	token := c.Query("token")
	valid, _ := V1VerifyToken(token)
	if valid {
		c.JSON(http.StatusOK, gin.H{
			"msg": "token valid",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "token invalid",
		})
	}

}

func V1GenerateToken(email string) string {
	log.Println("Generating token for user: ", email, " ...")
	signingMethod := jwt.SigningMethodHS256 // HS256 is an instance of HMAC
	claims := jwt.MapClaims{
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 7).Unix(),
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
		userPass := info[6:]
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
	if email == "chen-gz@look.com" && password == "Connie" {
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
