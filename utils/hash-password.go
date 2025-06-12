package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) string {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(hashPassword)

}
