package ports

import (
	"time"

	"github.com/google/uuid"
)

const maxPicturesCount = 10

type Content int

const (
	TextContent Content = iota
	PictureContent
	VideoContent
)

type SpecType uint8

const (
	UserSpecType SpecType = iota
	SpaceSpecType
)

type SpecAmount uint64
type MediaUrl string

// TODO convert to table - struct to keep metadata about tag and flair
type Spectrum uint64
type Flair string

type PostBody struct {
	QuoteId  uuid.UUID  `json:"title" gorm:""`
	Text     []string   `json:"title" gorm:""`
	Pictures []MediaUrl `json:"title" gorm:""`
	Video    MediaUrl   `json:"title" gorm:""`
	Order    []Content  `json:"title" gorm:""`
}

type Post struct {
	// TODO: change tags
	Id       uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;not null;"`
	UserId   uuid.UUID    `json:"user_id" gorm:"type:uuid;foreign_key;not null;"`
	Comments []*Post      `json:"comments" gorm:"foreignKey:Id;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags     []*uuid.UUID `json:"tags" gorm:"foreignKey:PostId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // TODO: recalculate tag list on Test update
	Spectrum Spectrum     `json:"spectrum" gorm:""`
	SpecType SpecType     `json:"spec_type" gorm:""`

	UrlKey string `json:"url_key" gorm:""`

	Flair   Flair `json:"flair" gorm:""`
	IsNSFW  bool  `json:"is_nsfw" gorm:""`
	Spoiler bool  `json:"spoiler" gorm:""`

	Title string   `json:"title" gorm:""`
	Body  PostBody `json:"body" gorm:""` // TODO: inline data in gorm

	TotalSpecAmount SpecAmount `json:"total_spec_amount" gorm:""`
	StarSpecAmount  SpecAmount `json:"star_spec_amount" gorm:""`
	OtherSpecAmount SpecAmount `json:"other_spec_amount" gorm:""`

	DonatedTokenQty map[string]string `json:"donated_token_qty" gorm:""` // presentation only - calculation delegated to other services

	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	IsDeleted bool      `json:"is_deleted" gorm:"not null"`
}
