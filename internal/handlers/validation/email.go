package validation

import (
	"errors"
	"regexp"
	"strings"
)

var (
	InvalidEmail = errors.New("email should end with @rutube.ru")
)

func Email(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	domain := strings.Split(email, "@")[1]

	if !emailRegex.MatchString(email) && !strings.HasSuffix(domain, "@rutube.ru") {
		return InvalidEmail
	}
	return nil
}
