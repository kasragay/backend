package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
	"time"
)

const postCaller = packageCaller + ".Post"

type Post struct {
	logger *utils.Logger
	cache  ports.CacheRepo
	rel    ports.RelationalRepo
	mongo  ports.MongoRepo
	s3     ports.S3Repo
}

func NewPostService(
	logger *utils.Logger,
	cache ports.CacheRepo,
	rel ports.RelationalRepo,
	mongo ports.MongoRepo,
	s3 ports.S3Repo,
) ports.PostService {
	return &Post{
		logger: logger,
		cache:  cache,
		rel:    rel,
		mongo:  mongo,
		s3:     s3,
	}
}

// TODO: implement service methods
func (s *Post) PostGet(ctx context.Context, req *ports.PostGetRequest) (resp *ports.Post, err error) {
	defer func() { err = utils.FuncPipe(postCaller+".PostGet", err) }()
	post, isDeleted, err := s.rel.GetPostById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, utils.PostDeletedResponse.Clone()
	}
	return &post, nil
}

func (s *Post) PostPost(ctx context.Context, req *ports.PostPostRequest) (resp *ports.PostPostResponse, err error) {
	defer func() { err = utils.FuncPipe(postCaller+".PostPost", err) }()

	post := req.ToPost()
	id, err := s.rel.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}
	post.Id = id
	return &ports.PostPostResponse{
		Id:        id,
		CreatedAt: post.CreatedAt,
	}, nil
}
func (s *Post) PostPut(ctx context.Context, req *ports.PostPutRequest) (err error) {
	defer func() { err = utils.FuncPipe(postCaller+".PostPut", err) }()

	err, done := s.checkForInvalidState(ctx, req.Id, err)
	if done {
		return err
	}

	updates := s.setUpdatedFields(req)

	err = s.rel.UpdatePostById(ctx, req.Id, updates)
	return err
}

func (s *Post) checkForInvalidState(ctx context.Context, id uuid.UUID, err error) (error, bool) {
	_, isDeleted, err := s.rel.GetPostById(ctx, id)
	if err != nil {
		return err, true
	}
	if isDeleted {
		return utils.PostDeletedResponse.Clone(), true
	}
	return nil, false
}

func (s *Post) setUpdatedFields(req *ports.PostPutRequest) map[string]interface{} {
	updates := make(map[string]interface{})

	// Top-level fields
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Flair != nil {
		updates["flair"] = *req.Flair
	}
	if req.IsNSFW != nil {
		updates["is_nsfw"] = *req.IsNSFW
	}
	if req.Spoiler != nil {
		updates["spoiler"] = *req.Spoiler
	}

	// Handle nested PostBody
	if req.Body != nil {
		bodyUpdates := make(map[string]interface{})
		if req.Body.QuoteId != nil {
			bodyUpdates["quote_id"] = *req.Body.QuoteId
		}
		if req.Body.Text != nil {
			bodyUpdates["text"] = req.Body.Text
		}
		if req.Body.Video != nil {
			bodyUpdates["video"] = *req.Body.Video
		}
		updates["post_body"] = bodyUpdates
	}

	updates["updated_at"] = time.Now().UTC()
	return updates
}

func (s *Post) PostDelete(ctx context.Context, id uuid.UUID) (err error) {
	defer func() { err = utils.FuncPipe(postCaller+".PostDelete", err) }()

	err, done := s.checkForInvalidState(ctx, id, err)
	if done {
		return err
	}

	err = s.rel.DeletePostById(ctx, id)

	return err
}
