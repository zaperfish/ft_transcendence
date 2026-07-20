package user

import (
    // Std
	"errors"
	"regexp"
	"strings"

    // External
	"github.com/go-ozzo/ozzo-validation/v4"
	// "github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	minUsernameLength = 3
	maxUsernameLength = 20
	minPasswordLength = 8
	maxPasswordLength = 128
)

func ValidUserName(name string) error {
	return validation.Validate(name, 
			validation.Required,
			validation.Length(minUsernameLength, maxUsernameLength),
			validation.Match(regexp.MustCompile(`^\S+$`)).Error("must not contain spaces"),
		)
}

var emailRegex = `^([A-Za-z0-9_%+-]+(?:\.[A-Za-z0-9_%+-]+)*)@(?:[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?\.)+[A-Za-z]{2,63}$`

func ValidUserEmail(email string) error {
	return validation.Validate(email, 
			validation.Required,
			makeRule(noConsecutiveDots),
			validation.Match(regexp.MustCompile(emailRegex)).Error("must be a valid email address"),
		)
}

func ValidUserPassword(password string) error {
	return validation.Validate(password, 
			validation.Required,
			validation.Length(minPasswordLength, maxPasswordLength),
			validation.Match(regexp.MustCompile(`[a-z]`)).Error("must contain at least one lower case character"),
			validation.Match(regexp.MustCompile(`[A-Z]`)).Error("must contain at least one upper case character"),
			validation.Match(regexp.MustCompile(`[0-9]`)).Error("must contain at least one digit"),
		)
}

func noConsecutiveDots(s string) error {
	if strings.Contains(s, `..`) {
		return errors.New("must be a valid email address")
	}
	return nil
}

func makeRule(test func(s string) error) validation.Rule {
	return validation.By(func(val any) error {
		s, _ := val.(string)
		return test(s)
	})
}
