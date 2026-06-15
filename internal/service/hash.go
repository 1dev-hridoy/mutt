package service

import (
	"github.com/dishan1223/mutt/consts"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) (string, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(p), consts.HASH_COST)
	if err != nil {
		return "", err
	}

	return string(hp), nil
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
