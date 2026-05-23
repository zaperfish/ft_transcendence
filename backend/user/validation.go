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
	maxUsernameLength = 30
	minPasswordLength = 4
	maxPasswordLength = 128
)

// func ValidateUser(u User) error {
// 	return validation.ValidateStruct(&u,
// 			validation.Field(&u.Name, makeRule(validUserName)),
// 			validation.Field(&u.Email, makeRule(validUserEmail)),
// 			validation.Field(&u.Password, makeRule(validUserPassword)),
// 		)
// }

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
		)
}

func makeRule(test func(s string) error) validation.Rule {
	return validation.By(func(val any) error {
		s, _ := val.(string)
		return test(s)
	})
}
