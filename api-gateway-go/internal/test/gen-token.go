package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("my-secret-key")

func GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(secret)
}

func main() {
	token, err := GenerateToken()
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
}
