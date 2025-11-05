package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Use the same secret key as your JWT middleware
	secretKey := []byte("supersecret")

	// Define your claims
	claims := jwt.MapClaims{
		"sub":  "user123",                              // subject (e.g. user ID)
		"role": "tester",                               // custom claim
		"exp":  time.Now().Add(5 * time.Minute).Unix(), // expiry time
		"iat":  time.Now().Unix(),                      // issued at
	}

	// Create token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Your JWT token:")
	fmt.Println(tokenString)
}
