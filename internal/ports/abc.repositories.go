package ports

import (
	"context"
	"image"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type RelationalRepo interface {
	UserExists(ctx context.Context, req *AuthCheckPostRequest) (resp *AuthCheckPostResponse, isDeleted bool, err error)
	UserExistsById(ctx context.Context, id uuid.UUID, userType UserType) (exists bool, isDeleted bool, err error)
	CreateUser(ctx context.Context, req *AuthSignupPostRequest, forceId ...uuid.UUID) (resp *User, err error)
	AuthUserPassword(ctx context.Context, req *AuthSigninPasswordPostRequest) (resp *User, err error)
	GetUserById(ctx context.Context, id uuid.UUID, userType UserType) (user UserModel, isDeleted bool, err error)
	GetUserByUsername(ctx context.Context, username string, userType UserType) (user UserModel, isDeleted bool, err error)
	UpdateUserPasswordById(ctx context.Context, id uuid.UUID, userType UserType, password string) (err error)
	UpdateUserPasswordByUsername(ctx context.Context, username string, userType UserType, password string) (err error)
	UpdateUserProfileById(ctx context.Context, id uuid.UUID, username, name, avatar string, userType UserType) (err error)
	DeleteUserById(ctx context.Context, id uuid.UUID, userType UserType) (err error)
	UpdateUserPhoneById(ctx context.Context, id uuid.UUID, userType UserType, phoneNumber string) (err error)
	UpdateUserEmailById(ctx context.Context, id uuid.UUID, userType UserType, email string) (err error)
	CheckUserEmailLimit(ctx context.Context, email string) (err error)
	CheckUserPhoneLimit(ctx context.Context, phoneNumber string) (err error)

	CreatePost(ctx context.Context, req *Post, forceId ...uuid.UUID) (postId uuid.UUID, err error)
	GetPostById(ctx context.Context, id uuid.UUID) (post Post, isDeleted bool, err error)
	UpdatePostById(ctx context.Context, postId uuid.UUID, fields map[string]any) (err error)
	DeletePostById(ctx context.Context, id uuid.UUID) (err error)

	Close() error

	DropTables() error
	AddFirstUsers() error // TODO_DEL
	AutoMigrate() error
}

type CacheRepo interface {
	AddJwtToBlacklist(ctx context.Context, token string, expire time.Duration) (err error)
	IsJwtInBlacklist(ctx context.Context, token string) (isIn bool, err error)

	SetOtpToken(ctx context.Context, identity, token string, otpType OtpType, expire time.Duration, userType UserType) (err error)
	GetOtpToken(ctx context.Context, identity string, otpType OtpType, userType UserType) (token string, err error)
	DeleteOtpToken(ctx context.Context, identity string, otpType OtpType, userType UserType) (err error)

	SetOtpKey(ctx context.Context, identity, key string, expire time.Duration, userType UserType) (err error)
	GetOtpKey(ctx context.Context, identity string, userType UserType) (key string, err error)
	DeleteOtpKey(ctx context.Context, identity string, userType UserType) (err error)

	Close() error
}

type S3Repo interface {
	GetAvatar(ctx context.Context, userId uuid.UUID, userType UserType) (avatar *minio.Object, err error)
	UploadAvatar(ctx context.Context, userId uuid.UUID, userType UserType, img *image.Image) (err error)
	DeleteAvatar(ctx context.Context, userId uuid.UUID, userType UserType) (err error)
}

type MongoRepo interface {
	Close() error
}
