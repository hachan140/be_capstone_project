package utils

import (
	"be-capstone-project/src/internal/core/events"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"regexp"
)

func ValidateAutoTestPhonePattern(phoneNumber string, autoPhonePattern string) bool {
	matched, _ := regexp.Match(autoPhonePattern, []byte(phoneNumber))
	return matched
}

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
