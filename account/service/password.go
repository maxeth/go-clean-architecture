package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// hashes a plain password and returns the hash on success.
func HashPassword(password string) (string, error) {
	hpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash pw %w", err)
	}

	return string(hpw), nil
}

// compares hash to given password and returns nil on success.
func ComparePassword(hPw, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hPw), []byte(pw))
}

// func hashPassword(pw string) (string, error) {
// 	salt := make([]byte, 32)
// 	_, err := rand.Read(salt)

// 	if err != nil {
// 		return "", err
// 	}

// 	sHash, err := scrypt.Key([]byte(pw), salt, 32768, 8, 1, 32)
// 	if err != nil {
// 		return "", err
// 	}

// 	hashedPw := fmt.Sprintf("%s.%s", hex.EncodeToString(sHash), hex.EncodeToString(salt))

// 	return hashedPw, nil
// }
