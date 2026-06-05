package user

import (
	"net/http"
)

// register

type CreateInput struct {
    Body CreateUserDTO
}

type CreateUserDTO struct {
    Name string     `json:"name" maxLength:"30" example:"Max" doc:"username"`
    Email string    `json:"email" example:"max@email.com" doc:"email address"`
    Password string `json:"password" example:"secret" doc:"password"`
}

// login

type UserLoginDTO struct {
    Name string     `json:"name" example:"Max"`
    Password string `json:"password" example:"secret"`
}

type LoginUserInput struct {
    Body UserLoginDTO
}

type LoginUserOutput struct {
	SetCookie http.Cookie 		`header:"Set-Cookie"`
    Body UserSummaryDTO
}

// logout

type LogoutUserOutput struct {
	SetCookie http.Cookie 		`header:"Set-Cookie"`
}

// get

type GetUserInput struct {
	ID uint	`path:"id" doc:"User ID" example:"1"`
}

// get list

type UserFilter struct {
	Page     int
	PageSize int
}

type GetUsersInput struct {
	Page     int `query:"page" minimum:"1" default:"1" doc:"Filter by page"`
	PageSize int `query:"page_size" minimum:"1" default:"10" doc:"Page size"`
}

// patch

type PatchUserInput struct {
	ID uint	`path:"id" doc:"User ID" example:"1"`
	Body PatchUserDTO
}

type PatchUserDTO struct {
    Name *string     `json:"name,omitempty" maxLength:"30" example:"Max" doc:"username"`
    Email *string    `json:"email,omitempty" example:"max@email.com" doc:"email address"`
}

// patch password

type PatchPasswordInput struct {
	ID uint	`path:"id" doc:"User ID" example:"1"`
	Body PatchPasswordDTO
}

type PatchPasswordDTO struct {
    NewPassword string `json:"newpassword" example:"newsecret" doc:"new password"`
    ConfirmPassword string `json:"confirm_password" example:"newsecret" doc:"confirm password"`
    CurrentPassword string `json:"current_password" example:"secret" doc:"current password"`
}

// delete

type DeleteUserInput struct {
	ID uint	`path:"id" doc:"User ID" example:"1"`
}
