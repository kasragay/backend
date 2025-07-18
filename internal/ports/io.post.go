package ports

import (
	"github.com/google/uuid"
	"time"
)

type PostGetRequest struct {
	Id uuid.UUID `json:"id" validate:"required,uuid4"`
}

type PostPutRequest struct {
	Id      uuid.UUID    `json:"id" validate:"required,uuid4"`
	Title   *string      `json:"title,omitempty"`
	Body    *PostPutBody `json:"body,omitempty"`
	Flair   *Flair       `json:"flair,omitempty" `
	IsNSFW  *bool        `json:"is_nsfw,omitempty" `
	Spoiler *bool        `json:"spoiler,omitempty" `
}

type PostPutBody struct {
	QuoteId  *uuid.UUID  `json:"quote_id,omitempty"`
	Text     []*string   `json:"text,omitempty" validate:"required_without_all= Pictures Video"`
	Pictures []*MediaUrl `json:"pictures,omitempty" validate:"required_without_all= Text Video"`
	Video    *MediaUrl   `json:"video,omitempty" validate:"required_without_all= Text Pictures"`
	Order    []*Content  `json:"order,omitempty" validate:"dive,required"`
}

// TODO: validate the string contents ( trim start & end )
type PostPostRequest struct {
	UserId   uuid.UUID    `json:"user_id" validate:"required,uuid4"`
	Title    string       `json:"title" validate:"required"`
	Body     PostPostBody `json:"body"`
	Flair    Flair        `json:"flair" `
	IsNSFW   bool         `json:"is_nsfw" `
	Spoiler  bool         `json:"spoiler" `
	Spectrum Spectrum     `json:"spectrum"`
	SpecType SpecType     `json:"spec_type"`
}

// TODO: change the media url arrays to 1-n relationships
type PostPostBody struct {
	QuoteId  uuid.UUID  `json:"quote_id,omitempty"`
	Text     []string   `json:"text,omitempty" validate:"required_without_all= Pictures Video Order"`
	Pictures []MediaUrl `json:"pictures,omitempty" validate:"required_without_all= Text Video Order"`
	Video    MediaUrl   `json:"video,omitempty" validate:"required_without_all= Text Pictures Order"`
	Order    []Content  `json:"order,omitempty" validate:"dive,required"`
}

func (req *PostPostRequest) ToPost() *Post {
	now := time.Now().UTC()

	post := &Post{
		UserId:    req.UserId,
		Spectrum:  req.Spectrum,
		SpecType:  req.SpecType,
		Flair:     req.Flair,
		IsNSFW:    req.IsNSFW,
		Spoiler:   req.Spoiler,
		Title:     req.Title,
		UpdatedAt: now,
		CreatedAt: now,
		IsDeleted: false,
	}

	post.PostBody = PostBody{
		QuoteId:  req.Body.QuoteId,
		Text:     req.Body.Text,
		Pictures: req.Body.Pictures,
		Video:    req.Body.Video,
		Order:    req.Body.Order,
	}

	return post
}

type PostPutResponse struct {
}

type PostPostResponse struct {
	Id        uuid.UUID `json:"id" validate:"required,uuid4"`
	CreatedAt time.Time `json:"created_at"`
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
