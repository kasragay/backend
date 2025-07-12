package services

import (
	"context"

	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
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
	// defer func() { err = utils.FuncPipe(postCaller+".PostGet", err) }()
	// user, isDeleted, err := s.rel.GetPostById(ctx, req.Id, req.PostType)
	// if err != nil {
	// 	return nil, err
	// }
	// if isDeleted {
	// 	return nil, utils.PostDeletedResponse.Clone()
	// }
	// return user.ToPost(), nil
	return
}

func (s *Post) PostPost(ctx context.Context, req *ports.PostPostRequest) (resp *ports.PostPostResponse, err error) {
	return
}
func (s *Post) PostPut(ctx context.Context, req *ports.PostPutRequest) (resp *ports.PostPutResponse, err error) {
	// defer func() { err = utils.FuncPipe(postCaller+".PostPut", err) }()

	// img, err := ports.AvatarValidator(req.Avatar)
	// if err != nil {
	// 	return nil, err
	// }
	// if img == nil {
	// 	if err := s.s3.DeleteAvatar(ctx, req.Id, req.PostType); err != nil {
	// 		return nil, err
	// 	}
	// } else {
	// 	if err := s.s3.UploadAvatar(ctx, req.Id, req.PostType, img); err != nil {
	// 		return nil, err
	// 	}
	// }
	// err = s.rel.UpdatePostProfileById(ctx, req.Id, req.Postname, req.Name, req.Avatar, req.PostType)
	// if err != nil {
	// 	return nil, err
	// }
	// return &ports.PostPostPutResponse{
	// 	Avatar: ports.GetAvatarUrl(req.Id, req.PostType),
	// }, nil
	return
}

func (s *Post) PostDelete(ctx context.Context, req *ports.PostDeleteRequest, tokenCheck bool) (err error) {
	// defer func() { err = utils.FuncPipe(postCaller+".PostDelete", err) }()
	// user, isDeleted, err := s.rel.GetPostById(ctx, req.Id, req.PostType)
	// if err != nil {
	// 	return err
	// }
	// if isDeleted {
	// 	return utils.PostDeletedResponse.Clone()
	// }
	// if user.GetHasAvatar() {
	// 	if err = s.s3.DeleteAvatar(ctx, req.Id, req.PostType); err != nil {
	// 		return err
	// 	}
	// }
	// if !tokenCheck {
	// 	return s.rel.DeletePostById(ctx, req.Id, req.PostType)
	// }
	// cToken, err := s.cache.GetOtpToken(ctx, user.GetPhoneNumber(), ports.DeleteAccountOtpType, req.PostType)
	// if err != nil {
	// 	return err
	// }
	// if cToken != req.Token {
	// 	return utils.TokenIncorrectResponse.Clone()
	// }
	// if err = s.cache.DeleteOtpToken(ctx, user.GetPhoneNumber(), ports.DeleteAccountOtpType, req.PostType); err != nil {
	// 	return err
	// }
	// return s.rel.DeletePostById(ctx, req.Id, req.PostType)
	return
}
