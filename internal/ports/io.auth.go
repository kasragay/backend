package ports

import (
	"github.com/google/uuid"
)

type JwtType string

const (
	AccessJwtType  JwtType = "access"
	RefreshJwtType JwtType = "refresh"
)

var JwtTypes = []JwtType{AccessJwtType, RefreshJwtType}

type UserType string

const (
	AdminUserType  UserType = "admin"
	ClientUserType UserType = "client"
)

var AllUserTypes = []UserType{AdminUserType, ClientUserType}

type AuthCheckPostRequest struct {
	Username string   `json:"username" validate:"required,usernameValidator"`
	UserType UserType `json:"user_type" validate:"required,userTypeValidator"`
}

type AuthCheckPostResponse struct {
	Exists         bool `json:"exists"`
	Deleted        bool `json:"deleted"`
	HasEmail       bool `json:"has_email"`
	HasPhoneNumber bool `json:"has_phone_number"`
	HasPassword    bool `json:"has_password"`
}

type OtpType string

const (
	AdminSignupKeyOtpType OtpType = "admin-signup-key" // only used for mailcom

	SignupOtpType         OtpType = "sign-up"
	SigninOtpType         OtpType = "sign-in"
	ChangePasswordOtpType OtpType = "change-password"
	ChangePhoneOtpType    OtpType = "change-phone"
	DeleteAccountOtpType  OtpType = "delete-account"
	ChangeEmailOtpType    OtpType = "change-recovery-email"
)

var OtpTypes = []OtpType{
	SignupOtpType,
	SigninOtpType,
	ChangePasswordOtpType,
	ChangePhoneOtpType,
	DeleteAccountOtpType,
	ChangeEmailOtpType,
}

type AuthMethodOtpGetRequest struct {
	Username    string   `json:"username" validate:"required,usernameValidator"`
	SendToEmail bool     `json:"send_to_email"`
	SendToPhone bool     `json:"send_to_phone"`
	Email       string   `json:"email" validate:"emailValidator"`
	PhoneNumber string   `json:"phone_number" validate:"phoneValidator"`
	OtpType     OtpType  `json:"otp_type" validate:"required,otpTypeValidator"`
	UserType    UserType `json:"user_type" validate:"required,userTypeValidator"`
}

type AuthMethodOtpGetResponse struct {
	MaskedPhone string `json:"masked_phone"`
	MaskedEmail string `json:"masked_email"`
}

type TmpAuthMethodOtpGetResponse struct {
	Token string `json:"token" validate:"required,otpTypeValidator"`
}

type AuthSignupKeyGetRequest struct {
	Email       string   `json:"email" validate:"emailValidator"`
	PhoneNumber string   `json:"phone_number" validate:"phoneValidator"`
	UserType    UserType `json:"user_type" validate:"required,nonClientUserTypeValidator"`
}

type TmpAuthSignupKeyGetResponse struct {
	Key string `json:"key" validate:"required,otpKeyValidator"`
}

type AuthSignupPostRequest struct {
	Username    string   `json:"username" validate:"required,usernameValidator"`
	Name        string   `json:"name" validate:"required,nameValidator"`
	Avatar      string   `json:"avatar,omitempty" validate:"avatarValidator"`
	PhoneNumber string   `json:"phone_number" validate:"phoneValidator"`
	Email       string   `json:"email" validate:"emailValidator"`
	UserType    UserType `json:"user_type" validate:"required,userTypeValidator"`
	Password    string   `json:"password" validate:"passwordValidator"`
	Token       string   `json:"token" validate:"required,otpTokenValidator"`
	Key         string   `json:"key" validate:"otpKeyValidator"`
}

type AuthSigninPostRequest struct {
	Username    string   `json:"username" validate:"required,usernameValidator"`
	SentToEmail bool     `json:"sent_to_email"`
	SentToPhone bool     `json:"sent_to_phone"`
	UserType    UserType `json:"user_type" validate:"required,userTypeValidator"`
	Token       string   `json:"token" validate:"required,otpTokenValidator"`
}

type AuthSigninPasswordPostRequest struct {
	Username string   `json:"username" validate:"required,usernameValidator"`
	UserType UserType `json:"user_type" validate:"required,userTypeValidator"`
	Password string   `json:"password" validate:"required,passwordValidator"`
}

type AuthResetPasswordPostRequest struct {
	Username    string   `json:"username" validate:"required,usernameValidator"`
	SentToEmail bool     `json:"sent_to_email"`
	SentToPhone bool     `json:"sent_to_phone"`
	Password    string   `json:"password" validate:"required,passwordValidator"`
	UserType    UserType `json:"user_type" validate:"required,userTypeValidator"`
	Token       string   `json:"token" validate:"required,otpTokenValidator"`
}

type AuthResetPhonePostRequest struct {
	Id          uuid.UUID `json:"id" validate:"required,uuid4"`
	PhoneNumber string    `json:"phone_number" validate:"required,phoneValidator"`
	UserType    UserType  `json:"user_type" validate:"required,userTypeValidator"`
	Token       string    `json:"token" validate:"required,otpTokenValidator"`
}

type AuthResetEmailPostRequest struct {
	Id       uuid.UUID `json:"id" validate:"required,uuid4"`
	Email    string    `json:"email" validate:"required,emailValidator"`
	UserType UserType  `json:"user_type" validate:"required,userTypeValidator"`
	Token    string    `json:"token" validate:"required,otpTokenValidator"`
}

type S3AvatarsGetRequest struct {
	UserType   UserType `json:"user_type" validate:"required,userTypeValidator"`
	ObjectName string   `json:"object_name" validate:"required,avatarObjectNameValidator"`
}

type Jwt struct {
	AccessToken    string `json:"access_token" validate:"required,jwt"`
	RefreshToken   string `json:"refresh_token" validate:"required,jwt"`
	AccessExpires  int    `json:"access_expires" validate:"required"`
	RefreshExpires int    `json:"refresh_expires" validate:"required"`
}

type Login struct {
	User *User `json:"user" validate:"required"`
	Jwt  *Jwt  `json:"jwt" validate:"required"`
}
