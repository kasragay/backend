package services

import (
	"context"

	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
)

const userCaller = packageCaller + ".User"

type User struct {
	logger *utils.Logger
	cache  ports.CacheRepo
	rel    ports.RelationalRepo
	mongo  ports.MongoRepo
	s3     ports.S3Repo
}

func NewUserService(
	logger *utils.Logger,
	cache ports.CacheRepo,
	rel ports.RelationalRepo,
	mongo ports.MongoRepo,
	s3 ports.S3Repo,
) ports.UserService {
	return &User{
		logger: logger,
		cache:  cache,
		rel:    rel,
		mongo:  mongo,
		s3:     s3,
	}
}

func (s *User) UserGet(ctx context.Context, req *ports.UserUserGetRequest) (resp *ports.User, err error) {
	defer func() { err = utils.FuncPipe(userCaller+".UserGet", err) }()
	user, isDeleted, err := s.rel.GetUserById(ctx, req.Id, req.UserType)
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, utils.UserDeletedResponse.Clone()
	}
	return user.ToUser(), nil
}

func (s *User) UserPut(ctx context.Context, req *ports.UserUserPutRequest) (resp *ports.UserUserPutResponse, err error) {
	defer func() { err = utils.FuncPipe(userCaller+".UserPut", err) }()

	img, err := ports.AvatarValidator(req.Avatar)
	if err != nil {
		return nil, err
	}
	if img == nil {
		if err := s.s3.DeleteAvatar(ctx, req.Id, req.UserType); err != nil {
			return nil, err
		}
	} else {
		if err := s.s3.UploadAvatar(ctx, req.Id, req.UserType, img); err != nil {
			return nil, err
		}
	}
	err = s.rel.UpdateUserProfileById(ctx, req.Id, req.Username, req.Name, req.Avatar, req.UserType)
	if err != nil {
		return nil, err
	}
	return &ports.UserUserPutResponse{
		Avatar: ports.GetAvatarUrl(req.Id, req.UserType),
	}, nil
}

func (s *User) UserDelete(ctx context.Context, req *ports.UserUserDeleteRequest, tokenCheck bool) (err error) {
	defer func() { err = utils.FuncPipe(userCaller+".UserDelete", err) }()
	user, isDeleted, err := s.rel.GetUserById(ctx, req.Id, req.UserType)
	if err != nil {
		return err
	}
	if isDeleted {
		return utils.UserDeletedResponse.Clone()
	}
	if user.GetHasAvatar() {
		if err = s.s3.DeleteAvatar(ctx, req.Id, req.UserType); err != nil {
			return err
		}
	}
	if !tokenCheck {
		return s.rel.DeleteUserById(ctx, req.Id, req.UserType)
	}
	cToken, err := s.cache.GetOtpToken(ctx, user.GetPhoneNumber(), ports.DeleteAccountOtpType, req.UserType)
	if err != nil {
		return err
	}
	if cToken != req.Token {
		return utils.TokenIncorrectResponse.Clone()
	}
	if err = s.cache.DeleteOtpToken(ctx, user.GetPhoneNumber(), ports.DeleteAccountOtpType, req.UserType); err != nil {
		return err
	}
	return s.rel.DeleteUserById(ctx, req.Id, req.UserType)
}
