package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secretKey := []byte("supersecret")

	claims := jwt.MapClaims{
		"sub":  "user123",                              
		"role": "tester",                               
		"exp":  time.Now().Add(5 * time.Minute).Unix(), 
		"iat":  time.Now().Unix(),                      
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Your JWT token:")
	fmt.Println(tokenString)
}
