package user

import (
    // Std
	"regexp"

    // External
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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

func ValidUserEmail(email string) error {
	return validation.Validate(email, 
			validation.Required,
			is.Email,
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

func makeRule(test func(s string) error) validation.Rule {
	return validation.By(func(val any) error {
		s, _ := val.(string)
		return test(s)
	})
}
