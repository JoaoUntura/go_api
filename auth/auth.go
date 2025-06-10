package auth

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GerarToken(userId uint) (string, error) {

	var jwt_secret = []byte(os.Getenv("JWT_SECRET"))

	payload := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString(jwt_secret)

}

func ValidarToken(tokenString string) (bool, *uint) {

	var jwt_secret = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifica se o método de assinatura é o esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inválido")
		}
		return jwt_secret, nil
	})

	if err != nil || !token.Valid {
		return false, nil
	} else {
		if payload, ok := token.Claims.(jwt.MapClaims); ok {
			userid := payload["user_id"].(float64)
			userIdUint := uint(userid)
			return true, &userIdUint
		} else {
			return false, nil
		}
	}
}
