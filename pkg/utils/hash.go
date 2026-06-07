package utils

import (
	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, error) {

    passwordBytes := []byte(password)

    hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }

    return string(hash), nil

}
func CheckPasswordHash (password, hash string) bool {

    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil


}
