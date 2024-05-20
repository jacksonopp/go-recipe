package services

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const DEFAULT_TIMEOUT = 5 * time.Second

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

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

func genRandStr(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes), nil
}
