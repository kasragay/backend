package ports

import "github.com/google/uuid"

// TODO: implement Post put and post methods - discuss

type PostGetRequest struct {
	Id uuid.UUID `json:"id" validate:"required,uuid4"`
}

type PostPutRequest struct {
	Id    uuid.UUID `json:"id" validate:"required,uuid4"`
	Title string    `json:"title"`
	Body  PostBody  `json:"body"`
	// TODO: discuss tags then add them to the put request
	Flair   Flair `json:"flair" gorm:""`
	IsNSFW  bool  `json:"is_nsfw" gorm:""`
	Spoiler bool  `json:"spoiler" gorm:""`
}

type PostPostRequest struct {
	Id    uuid.UUID `json:"id" validate:"required,uuid4"`
	Title string    `json:"title"`
	Body  PostBody  `json:"body"`
	// TODO: discuss tags then add them to the put request
	Flair   Flair `json:"flair" gorm:""`
	IsNSFW  bool  `json:"is_nsfw" gorm:""`
	Spoiler bool  `json:"spoiler" gorm:""`
}

type PostPutResponse struct {
}

type PostPostResponse struct {
}

type PostDeleteRequest struct {
	Id    uuid.UUID `json:"id" validate:"required,uuid4"`
	Token string    `json:"token"`
}

// type Post struct {
// 	Id       uuid.UUID `json:"id" validate:"required,uuid4"`
// 	Postname string    `json:"username" validate:"required,usernameValidator"`
// 	Name     string    `json:"name" validate:"required,nameValidator"`
// 	Avatar   string    `json:"avatar"`
// 	PostType PostType  `json:"user_type" validate:"required,userTypeValidator"`
// }
