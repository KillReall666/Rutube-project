package validation

import (
	"errors"
	"regexp"
)

var (
	InvalidPhoneNumber = errors.New("phone number format must be +79991234567")
)

func PhoneNumber(phoneNumber string) error {
	regex := regexp.MustCompile(`^\+7\d{10}$`)
	if !regex.MatchString(phoneNumber) {
		return InvalidPhoneNumber
	}
	return nil
}
