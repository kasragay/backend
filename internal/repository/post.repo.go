package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
	"gorm.io/gorm"
)

func (s *Relational) CreatePost(ctx context.Context, req *ports.Post, forceId ...uuid.UUID) (id *uuid.UUID, err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".CreatePost", err) }()
	postId := uuid.New()
	err = s.client.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			if len(forceId) > 0 {
				postId = forceId[0]
			}
			req.Id = postId
			if err := tx.WithContext(ctx).Create(*req).Error; err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &postId, nil
}

func (s *Relational) GetPostById(ctx context.Context, id uuid.UUID) (post *ports.Post, isDeleted bool, err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".GetPostById", err) }()

	post = &ports.Post{}
	if err := s.client.WithContext(ctx).Where("id = ?", id).First(post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if post.IsDeleted {
		return nil, true, nil
	}
	return post, false, nil
}

func (s *Relational) UpdatePostById(ctx context.Context, postId uuid.UUID, fields map[string]any) (err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".UpdatePostById", err) }()

	result := s.client.Model(&ports.Post{}).
		Where("id = ?", postId).
		Select("*").
		Updates(fields)

	return result.Error
}

func (s *Relational) DeletePostById(ctx context.Context, id uuid.UUID) (err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".DeletePostById", err) }()
	result := s.client.WithContext(ctx).Delete(&ports.Post{}, id)
	return result.Error
}
