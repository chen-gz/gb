// package main
package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var secreteKey = []byte("bcb967bec859b86e96564992792636bb442548af35a2e3374cee7a0f92542c18")

func V1Login(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")

	if email == "test" && password == "test" {
		c.JSON(http.StatusOK, gin.H{
			"token": V1ReleaseToken(email),
			"msg":   "log in success",
			"name":  "test",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "log in failed",
		})
	}
}

func V1ReleaseToken(user string) string {
	signingMethod := jwt.SigningMethodHS256
	//signingMethod := jwt.SigningMethodHMAC{}
	claims := jwt.MapClaims{
		"user": user,
		"name": "test",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	token := jwt.NewWithClaims(signingMethod, claims)
	signedToken, err := token.SignedString(secreteKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated token: ", signedToken)
	//Verify token
	return signedToken

}

func V1VerifyToken(token string) (bool, interface{}) {
	parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secreteKey, nil
	})
	valid := parsedToken.Valid
	email := parsedToken.Claims.(jwt.MapClaims)["email"]
	return valid, email
}
