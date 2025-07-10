package ports

import (
	"github.com/google/uuid"
)

type UserUserGetRequest struct {
	Id       uuid.UUID `json:"id" validate:"required,uuid4"`
	UserType UserType  `json:"user_type" validate:"required,userTypeValidator"`
}

type UserUserPutRequest struct {
	Id       uuid.UUID `json:"id" validate:"required,uuid4"`
	Username string    `json:"username" validate:"required,usernameValidator"`
	Name     string    `json:"name" validate:"required,nameValidator"`
	Avatar   string    `json:"avatar" validate:"avatarValidator"`
	UserType UserType  `json:"user_type" validate:"required,userTypeValidator"`
}

type UserUserPutResponse struct {
	Avatar string `json:"avatar"`
}

type UserUserDeleteRequest struct {
	Id       uuid.UUID `json:"id" validate:"required,uuid4"`
	UserType UserType  `json:"user_type" validate:"required,userTypeValidator"`
	Token    string    `json:"token"`
}

type User struct {
	Id       uuid.UUID `json:"id" validate:"required,uuid4"`
	Username string    `json:"username" validate:"required,usernameValidator"`
	Name     string    `json:"name" validate:"required,nameValidator"`
	Avatar   string    `json:"avatar"`
	UserType UserType  `json:"user_type" validate:"required,userTypeValidator"`
}
