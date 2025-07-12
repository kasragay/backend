package ports

import "github.com/google/uuid"

type Tag struct {
	Name       string    `json:"id" gorm:"type:string;primary_key;not null;"`
	PostId     uuid.UUID `json:"post_id" gorm:"type:uuid;foreign_key;not null;"`
	TotalCount uint64    `json:"total_count"`
}
