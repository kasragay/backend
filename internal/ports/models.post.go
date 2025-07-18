package ports

import (
	"github.com/google/uuid"
	"time"
)

type PostModel interface {
	New(id uuid.UUID, userId uuid.UUID, spectrum Spectrum,
		specType SpecType, flair Flair, isNSFW bool, spoiler bool,
		title string, body PostBody) Post
}

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

type Spectrum uint64
type Flair string

type PostBody struct {
	QuoteId  uuid.UUID  `json:"quote_id,omitempty"`
	Text     []string   `json:"text,omitempty"`
	Pictures []MediaUrl `json:"pictures,omitempty" gorm:"type:json"`
	Video    MediaUrl   `json:"video,omitempty"`
	Order    []Content  `json:"order,omitempty" gorm:"type:json"`
}

type Post struct {
	// TODO: change tags
	Id       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;not null;"`
	UserId   uuid.UUID `json:"user_id" gorm:"type:uuid;foreign_key;not null;"`
	Comments []*Post   `json:"comments" gorm:"foreignKey:Id;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags     []*Tag    `json:"tags" gorm:"many2many:post_tags"` // TODO: recalculate tag list on Test update
	Spectrum Spectrum  `json:"spectrum" gorm:""`
	SpecType SpecType  `json:"spec_type" gorm:""`

	UrlKey string `json:"url_key" gorm:""`

	Flair   Flair `json:"flair" gorm:""`
	IsNSFW  bool  `json:"is_nsfw" gorm:""`
	Spoiler bool  `json:"spoiler" gorm:""`

	Title    string                `json:"title" gorm:""`
	PostBody `json:"body" gorm:""` // TODO: inline data in gorm

	TotalSpecAmount SpecAmount `json:"total_spec_amount" gorm:""`
	StarSpecAmount  SpecAmount `json:"star_spec_amount" gorm:""`
	OtherSpecAmount SpecAmount `json:"other_spec_amount" gorm:""`

	DonatedTokenQty map[string]string `json:"donated_token_qty" gorm:""` // presentation only - calculation delegated to other services

	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	IsDeleted bool      `json:"is_deleted" gorm:"not null"`
}

func NewPost(id uuid.UUID, userId uuid.UUID, spectrum Spectrum, specType SpecType, flair Flair, isNSFW bool, spoiler bool, title string, body PostBody) Post {
	return Post{}
}
