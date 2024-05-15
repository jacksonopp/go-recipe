package services

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password, salt string) (string, error) {
	passAndSalt := password + salt
	bytes, err := bcrypt.GenerateFromPassword([]byte(passAndSalt), 14)
	return string(bytes), err
}

func checkPasswordHash(password, salt, hash string) bool {
	passAndSalt := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passAndSalt))
	return err == nil
}

func genRandStr(length int) (string, error) {
	result := make([]byte, length)
	_, err := rand.Read(result)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		result[i] &= 0x7F                       // Ensure the byte is within printable ASCII range
		for result[i] < 33 || result[i] > 126 { // Exclude control characters and space
			_, err = rand.Read(result[i : i+1])
			if err != nil {
				return "", err
			}
			result[i] &= 0x7F
		}
	}

	return string(result), nil
}
