package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secretKey = []byte("secret-key")

func Generate(userID string) (string, int64, error) {
	expiresIn := time.Now().Add(24 * time.Hour).Unix()
	// Kenapa JWT?
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": userID,
			"exp":    expiresIn,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresIn, nil
}

func VerifyAndCheckExpiration(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				fmt.Println(time.Unix(int64(exp), 0))
				fmt.Println(time.Now())
				return false, fmt.Errorf("token is expired")
			}
			return true, nil
		}
	}
	return false, fmt.Errorf("invalid token claims")
}
