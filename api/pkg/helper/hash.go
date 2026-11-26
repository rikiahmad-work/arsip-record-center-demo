package helper

import (
	"golang.org/x/crypto/bcrypt"
)

// HashHelper menyediakan fungsi hashing password.
type HashHelper struct {
	cost int
}

// NewHashHelper membuat instance HashHelper dengan cost default (12).
func NewHashHelper() *HashHelper {
	return &HashHelper{cost: bcrypt.DefaultCost + 2} // cost = 12
}

// Hash membuat bcrypt hash dari password plain text.
func (h *HashHelper) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify membandingkan password plain text dengan hash.
// Mengembalikan true jika cocok.
func (h *HashHelper) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
