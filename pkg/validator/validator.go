package validator

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func ValidateEmail(email string) bool {
	email = strings.TrimSpace(email)
	return emailRegex.MatchString(email)
}

func ValidateRequired(value string) bool {
	return strings.TrimSpace(value) != ""
}

func ValidateMinLength(value string, min int) bool {
	return len(strings.TrimSpace(value)) >= min
}

func ValidateMaxLength(value string, max int) bool {
	return len(strings.TrimSpace(value)) <= max
}

