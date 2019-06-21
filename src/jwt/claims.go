package jwt

import "github.com/dgrijalva/jwt-go"

// ArithmeticCustomClaims 自定义声明
type ArithmeticCustomClaims struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`

	jwt.StandardClaims
}
