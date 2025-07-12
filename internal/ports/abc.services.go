package ports

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserService interface {
	UserGet(ctx context.Context, req *UserUserGetRequest) (resp *User, err error)
	UserPut(ctx context.Context, req *UserUserPutRequest) (resp *UserUserPutResponse, err error)
	UserDelete(ctx context.Context, req *UserUserDeleteRequest, tokenCheck bool) (err error)
}

type PostService interface {
	PostGet(ctx context.Context, req *PostGetRequest) (resp *Post, err error)
	PostPost(ctx context.Context, req *PostPostRequest) (resp *PostPostResponse, err error)
	PostPut(ctx context.Context, req *PostPutRequest) (resp *PostPutResponse, err error)
	PostDelete(ctx context.Context, req *PostDeleteRequest, tokenCheck bool) (err error)
}

type AuthService interface {
	CheckPost(ctx context.Context, req *AuthCheckPostRequest) (resp *AuthCheckPostResponse, err error)
	CheckById(ctx context.Context, id uuid.UUID, userType UserType) (exists bool, isDeleted bool, err error)
	TmpSignupKeyGet(ctx context.Context, req *AuthSignupKeyGetRequest) (resp *TmpAuthSignupKeyGetResponse, err error)
	SignupKeyGet(ctx context.Context, req *AuthSignupKeyGetRequest) (err error)
	TmpMethodOtpGet(ctx context.Context, req *AuthMethodOtpGetRequest) (resp *TmpAuthMethodOtpGetResponse, err error)
	MethodOtpGet(ctx context.Context, req *AuthMethodOtpGetRequest) (resp *AuthMethodOtpGetResponse, err error)
	SignupPost(ctx context.Context, req *AuthSignupPostRequest) (resp *Login, err error)
	SigninPost(ctx context.Context, req *AuthSigninPostRequest) (resp *Login, err error)
	SigninPasswordPost(ctx context.Context, req *AuthSigninPasswordPostRequest) (resp *Login, err error)
	LogoutPost(ctx context.Context, login *Login) (err error)
	RefreshPost(ctx context.Context, login *Login) (resp *Jwt, err error)
	ResetPasswordPost(ctx context.Context, req *AuthResetPasswordPostRequest) (resp *Login, err error)
	ResetPhonePost(ctx context.Context, req *AuthResetPhonePostRequest) (err error)
	ResetEmailPost(ctx context.Context, req *AuthResetEmailPostRequest) (err error)

	// Internal endpoints
	SendOtp(ctx context.Context, req *AuthMethodOtpGetRequest) (resp *AuthMethodOtpGetResponse, err error)
	SendKey(ctx context.Context, req *AuthSignupKeyGetRequest) (err error)

	GenerateToken(ctx context.Context, user *Login) (err error)
	ParseToken(token string, jwtType JwtType) (login *Login, err error)
}

type TelecomService interface {
	NoReplySend(ctx context.Context, dst []string, message string) (success []string, err error)
}

type MailcomService interface {
	NoReplySend(ctx context.Context, dst []string, subject, message string) (success []string, err error)
}

type RatelimiterService interface {
	Handler(router string, method Method, path string, token uint64, duration time.Duration, onIP, onID bool) fiber.Handler
}

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	OPTIONS Method = "OPTIONS"
	PATCH   Method = "PATCH"
)
