package app

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	customErros "social-network/api-gateway/internal/errors"
	"social-network/api-gateway/internal/logger"
)

type Claims struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

func JWTTokenVerify(r *http.Request) error {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		logger.Error("token is empty")
		return &customErros.JWTTokenEmpty{}
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwtKey, _ := os.LookupEnv("SECRET_KEY")
		return []byte(jwtKey), nil
	})

	if err != nil || !token.Valid {
		logger.Error("token is invalid")
		return &customErros.JWTTokenInvalid{}
	}

	r.Header.Set("login", claims.Login)
	r.Header.Set("name", claims.Name)

	return nil
}
