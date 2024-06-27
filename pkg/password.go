package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"log"
)

func HashPassword(input []byte) (string, error) {
	hasher := md5.New()

	_, err := hasher.Write(input)
	if err != nil {
		log.Println("error hashing password")
		return "", err
	}

	digest := hasher.Sum(nil)

	result := hex.EncodeToString(digest)

	return result, nil
}

func ValidatePassword(password, hash string) bool {
	hashedPassword, err := HashPassword([]byte(password))
	return err == nil && hash == hashedPassword
}
