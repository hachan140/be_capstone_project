package utils

import (
	"be-capstone-project/src/internal/core/events"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func EncodeEvent(event events.Event) ([]byte, error) {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return nil, err
	}

	hash := md5.New()
	_, err = hash.Write(payload)
	if err != nil {
		return nil, err
	}

	event.PayloadID = hex.EncodeToString(hash.Sum(nil))

	evenByte, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return evenByte, nil
}

func EncryptPassword(password string) (string, error) {
	hashesPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashesPassword), nil
}

func CheckHash(hpass string, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hpass), []byte(pass))
	if err != nil {
		return err
	}
	return nil
}

func IsEmail(emailToCheck string) bool {
	if strings.Contains(emailToCheck, "@") {
		return true
	}
	return false
}
