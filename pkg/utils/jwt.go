package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID  uint    `json:"user_id"`
    Role    string  `json:"role"`
    jwt.RegisteredClaims
}

var jwtSecretKey = []byte(getEnv("JWT_SECRET", "very_long_fallback_development_only_secret_key_123456"))

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists{
        return value
    }
    return fallback
}


func GenerateToken(userID uint, role string) (string, error) {
    claim := Claims{
        UserID: userID,
        Role: role,
        RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour))},
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

    tokenString, err := token.SignedString(jwtSecretKey)

    if err != nil {
        return "", err
    }

    return tokenString, nil

}

func VerifyToken(tokenString string) (*Claims, error) {
    claim := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claim, 
    func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jwtSecretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claim, nil
}
