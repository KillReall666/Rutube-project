package validation

import (
	"regexp"
	"strings"
)

func Email(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	domain := strings.Split(email, "@")[1]

	if emailRegex.MatchString(email) && strings.HasSuffix(domain, "@rutube.ru") {
		return true
	}
	return false
}
