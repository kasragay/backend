package ports

type Tag struct {
	Name       string  `json:"id" gorm:"type:string;primary_key;not null;"`
	Posts      []*Post `json:"post_id" gorm:"many2many:post_tags;"`
	TotalCount uint64  `json:"total_count"`
}
