package services

import (
	"context"
	"math/rand"

	"os"
	"time"

	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
)

const authCaller = packageCaller + ".Auth"

type Auth struct {
	logger        *utils.Logger
	cache         ports.CacheRepo
	rel           ports.RelationalRepo
	mongo         ports.MongoRepo
	s3            ports.S3Repo
	telecom       ports.TelecomService
	mailcom       ports.MailcomService
	jwtSK         []byte
	jwtAccessExp  time.Duration
	jwtRefreshExp time.Duration
	domain        string
	version       string
	supportEmail  string
}

func NewAuthService(
	logger *utils.Logger,
	cr ports.CacheRepo,
	rel ports.RelationalRepo,
	s3 ports.S3Repo,
	mongo ports.MongoRepo,
	telecom ports.TelecomService,
	mailcom ports.MailcomService,
) ports.AuthService {
	jwtSK := os.Getenv("JWT_SECRET_KEY")
	if jwtSK == "" {
		logger.Fatal(context.Background(), "JWT_SECRET_KEY is not set")
	}

	jwtAE, err := utils.GetenvAsMinuteDuration("JWT_ACCESS_EXP", 0, true)
	if err != nil {
		logger.Fatal(context.Background(), err.Error())
	}
	jwtRE, err := utils.GetenvAsMinuteDuration("JWT_REFRESH_EXP", 0, true)
	if err != nil {
		logger.Fatal(context.Background(), err.Error())
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		logger.Fatal(context.Background(), "DOMAIN is not set")
	}

	version := os.Getenv("VERSION")
	if version == "" {
		logger.Fatal(context.Background(), "VERSION is not set")
	}
	supportEmail := os.Getenv("SUPPORT_EMAIL")
	if supportEmail == "" {
		logger.Fatal(context.Background(), "SUPPORT_EMAIL is not set")
	}
	return &Auth{
		logger:        logger,
		cache:         cr,
		rel:           rel,
		mongo:         mongo,
		s3:            s3,
		telecom:       telecom,
		mailcom:       mailcom,
		jwtSK:         []byte(jwtSK),
		jwtAccessExp:  jwtAE,
		jwtRefreshExp: jwtRE,
		domain:        domain,
		version:       version,
		supportEmail:  supportEmail,
	}
}

func (s *Auth) CheckPost(ctx context.Context, req *ports.AuthCheckPostRequest) (resp *ports.AuthCheckPostResponse, err error) {
	defer func() { err = utils.FuncPipe(authCaller+".CheckPost", err) }()
	resp, _, err = s.rel.UserExists(ctx, req)
	return
}

func (s *Auth) CheckById(ctx context.Context, id uuid.UUID, userType ports.UserType) (exists, isDeleted bool, err error) {
	defer func() { err = utils.FuncPipe(authCaller+".CheckById", err) }()
	return s.rel.UserExistsById(ctx, id, userType)
}

func (s *Auth) TmpSignupKeyGet(ctx context.Context, req *ports.AuthSignupKeyGetRequest) (resp *ports.TmpAuthSignupKeyGetResponse, err error) {
	defer func() { err = utils.FuncPipe(authCaller+".TmpSignupKeyGet", err) }()
	var identity string
	if req.Email != "" {
		identity = req.Email
	} else {
		identity = req.PhoneNumber
	}
	key, err := s.cache.GetOtpKey(ctx, identity, req.UserType)
	if err != nil {
		return nil, err
	}
	return &ports.TmpAuthSignupKeyGetResponse{Key: key}, nil
}

func (s *Auth) SignupKeyGet(ctx context.Context, req *ports.AuthSignupKeyGetRequest) (err error) {
	defer func() { err = utils.FuncPipe(authCaller+".SignupKeyGet", err) }()
	return s.SendKey(ctx, req)
}

func (s *Auth) TmpMethodOtpGet(ctx context.Context, req *ports.AuthMethodOtpGetRequest) (resp *ports.TmpAuthMethodOtpGetResponse, err error) {
	defer func() { err = utils.FuncPipe(authCaller+".TmpMethodOtpGet", err) }()

	var token string
	var usr ports.UserModel
	usr, _, err = s.rel.GetUserByUsername(ctx, req.Username, req.UserType)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, utils.UsernameNotFoundResponse.Clone()
	}
	if req.SendToEmail {
		if usr.GetEmail() == "" {
			return nil, utils.UserHasNotSetEmailResponse.Clone().
				WithReason("username", req.Username)
		}
		token, err = s.cache.GetOtpToken(ctx, usr.GetEmail(), req.OtpType, req.UserType)
	} else {
		token, err = s.cache.GetOtpToken(ctx, usr.GetPhoneNumber(), req.OtpType, req.UserType)
	}
	if err != nil {
		return nil, err
	}
	return &ports.TmpAuthMethodOtpGetResponse{Token: token}, nil
}

func (s *Auth) MethodOtpGet(ctx context.Context, req *ports.AuthMethodOtpGetRequest) (resp *ports.AuthMethodOtpGetResponse, err error) {
	defer func() { err = utils.FuncPipe(authCaller+".MethodOtpGet", err) }()
	return s.SendOtp(ctx, req)
}

func (s *Auth) SignupPost(ctx context.Context, req *ports.AuthSignupPostRequest) (resp *ports.Login, err error) {
	defer func() { err = utils.FuncPipe(authCaller+".SignupAdminPost", err) }()
	var identity string
	if req.Email != "" {
		identity = req.Email
	} else {
		identity = req.PhoneNumber
	}
	if req.UserType == ports.AdminUserType {
		if req.Key == "" {
			return nil, utils.AuthSignupKeyIncorrectResponse.Clone()
		}
		cKey, err := s.cache.GetOtpKey(ctx, identity, req.UserType)
		if err != nil {
			return nil, err
		}
		if cKey != req.Key {
			return nil, utils.AuthSignupKeyIncorrectResponse.Clone()
		}
		if err = s.cache.DeleteOtpKey(ctx, identity, req.UserType); err != nil {
			return nil, err
		}
	}

	cToken, err := s.cache.GetOtpToken(ctx, identity, ports.SignupOtpType, req.UserType)
	if err != nil {
		return nil, err
	}
	if cToken != req.Token {
		return nil, utils.TokenIncorrectResponse.Clone()
	}
	if err = s.cache.DeleteOtpToken(ctx, identity, ports.SignupOtpType, req.UserType); err != nil {
		return nil, err
	}
	img, err := ports.AvatarValidator(req.Avatar)
	if err != nil {
		return nil, err
	}
	hasAvatar := img != nil
	resp = &ports.Login{}
	resp.User, err = s.rel.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if hasAvatar {
		if err = s.s3.UploadAvatar(ctx, resp.User.Id, req.UserType, img); err != nil {
			return nil, err
		}
	}
	err = s.GenerateToken(ctx, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Auth) SigninPost(ctx context.Context, req *ports.AuthSigninPostRequest) (resp *ports.Login, err error) {
	defer func() { err = utils.FuncPipe(authCaller+".SigninPost", err) }()
	usrModel, isDeleted, err := s.rel.GetUserByUsername(ctx, req.Username, req.UserType)
	if err != nil {
		return nil, err
	}
	if usrModel == nil {
		return nil, utils.UsernameNotFoundResponse.Clone().
			WithReason("username", req.Username)
	}
	if isDeleted {
		return nil, utils.UserDeletedResponse.Clone()
	}
	var identity string
	if req.SentToEmail {
		identity = usrModel.GetEmail()
	} else {
		identity = usrModel.GetPhoneNumber()
	}

	cToken, err := s.cache.GetOtpToken(ctx, identity, ports.SigninOtpType, req.UserType)
	if err != nil {
		return nil, err
	}
	if cToken != req.Token {
		return nil, utils.TokenIncorrectResponse.Clone()
	}
	if err = s.cache.DeleteOtpToken(ctx, identity, ports.SigninOtpType, req.UserType); err != nil {
		return nil, err
	}

	resp = &ports.Login{
		User: usrModel.ToUser(),
	}
	err = s.GenerateToken(ctx, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Auth) SigninPasswordPost(ctx context.Context, req *ports.AuthSigninPasswordPostRequest) (resp *ports.Login, err error) {
	defer func() { err = utils.FuncPipe(authCaller+".SigninPasswordPost", err) }()
	user, _, err := s.rel.GetUserByUsername(ctx, req.Username, req.UserType)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.UsernameNotFoundResponse.Clone().
			WithReason("username", req.Username)
	}
	ok := utils.VerifyPassword(user.GetPassword(), req.Password)
	if !ok {
		return nil, utils.PasswordIncorrectResponse.Clone()
	}
	resp = &ports.Login{
		User: user.ToUser(),
	}
	err = s.GenerateToken(ctx, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Auth) LogoutPost(ctx context.Context, login *ports.Login) (err error) {
	defer func() { err = utils.FuncPipe(authCaller+".LogoutPost", err) }()
	if err := s.cache.AddJwtToBlacklist(ctx, login.Jwt.AccessToken, s.jwtAccessExp); err != nil {
		return err
	}
	if err := s.cache.AddJwtToBlacklist(ctx, login.Jwt.RefreshToken, s.jwtRefreshExp); err != nil {
		return err
	}
	return
}

func (s *Auth) RefreshPost(ctx context.Context, login *ports.Login) (resp *ports.Jwt, err error) {
	defer func() { err = utils.FuncPipe(authCaller+".RefreshPost", err) }()
	if err = s.GenerateToken(ctx, login); err != nil {
		return nil, err
	}
	return &ports.Jwt{
		AccessToken:    login.Jwt.AccessToken,
		RefreshToken:   login.Jwt.RefreshToken,
		AccessExpires:  login.Jwt.AccessExpires,
		RefreshExpires: login.Jwt.RefreshExpires,
	}, nil
}

func (s *Auth) ResetPasswordPost(ctx context.Context, req *ports.AuthResetPasswordPostRequest) (resp *ports.Login, err error) {
	defer func() { err = utils.FuncPipe(authCaller+".ResetPasswordPost", err) }()

	user, isDeleted, err := s.rel.GetUserByUsername(ctx, req.Username, req.UserType)
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, utils.UserDeletedResponse.Clone()
	}
	if user == nil {
		return nil, utils.UsernameNotFoundResponse.Clone().
			WithReason("username", req.Username)
	}

	var identity string
	if req.SentToEmail {
		identity = user.GetEmail()
	} else {
		identity = user.GetPhoneNumber()
	}

	cToken, err := s.cache.GetOtpToken(ctx, identity, ports.ChangePasswordOtpType, req.UserType)
	if err != nil {
		return nil, err
	}
	if cToken != req.Token {
		return nil, utils.TokenIncorrectResponse.Clone()
	}
	if err = s.cache.DeleteOtpToken(ctx, identity, ports.ChangePasswordOtpType, req.UserType); err != nil {
		return nil, err
	}
	if err = s.rel.UpdateUserPasswordByUsername(ctx, req.Username, req.UserType, req.Password); err != nil {
		return nil, err
	}

	resp = &ports.Login{
		User: user.ToUser(),
	}
	err = s.GenerateToken(ctx, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Auth) ResetPhonePost(ctx context.Context, req *ports.AuthResetPhonePostRequest) (err error) {
	defer func() { err = utils.FuncPipe(authCaller+".ResetPhonePost", err) }()

	user, isDeleted, err := s.rel.GetUserById(ctx, req.Id, req.UserType)
	if err != nil {
		return err
	}
	if isDeleted {
		return utils.UserDeletedResponse.Clone()
	}
	if user == nil {
		return utils.UserNotFoundResponse.Clone().
			WithReason("id", req.Id)
	}
	prevPhone := user.GetPhoneNumber()
	cToken, err := s.cache.GetOtpToken(ctx, prevPhone, ports.ChangePhoneOtpType, req.UserType)
	if err != nil {
		return err
	}
	if cToken != req.Token {
		return utils.TokenIncorrectResponse.Clone()
	}
	if err = s.cache.DeleteOtpToken(ctx, prevPhone, ports.ChangePhoneOtpType, req.UserType); err != nil {
		return err
	}
	if err = s.rel.UpdateUserPhoneById(ctx, req.Id, req.UserType, req.PhoneNumber); err != nil {
		return err
	}
	return nil
}

func (s *Auth) ResetEmailPost(ctx context.Context, req *ports.AuthResetEmailPostRequest) (err error) {
	defer func() { err = utils.FuncPipe(authCaller+".ResetEmailPost", err) }()
	user, deleted, err := s.rel.GetUserById(ctx, req.Id, req.UserType)
	if err != nil {
		return err
	}
	if deleted {
		return utils.UserDeletedResponse.Clone()
	}
	if user == nil {
		return utils.UserNotFoundResponse.Clone().
			WithReason("id", req.Id)
	}
	prevEmail := user.GetEmail()
	cToken, err := s.cache.GetOtpToken(ctx, prevEmail, ports.ChangeEmailOtpType, req.UserType)
	if err != nil {
		return err
	}
	if cToken != req.Token {
		return utils.TokenIncorrectResponse.Clone()
	}
	if err = s.cache.DeleteOtpToken(ctx, prevEmail, ports.ChangeEmailOtpType, req.UserType); err != nil {
		return err
	}
	if err = s.rel.UpdateUserEmailById(ctx, req.Id, req.UserType, req.Email); err != nil {
		return err
	}
	return nil
}

// Internal endpoints

func generateRandomOTP() string {
	return fmt.Sprintf("%05d", utils.RandomInt(10000, 99999))
}

func generateRandomKey() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	result := make([]byte, 8)
	for i := range 8 {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func (s *Auth) SendOtp(ctx context.Context, req *ports.AuthMethodOtpGetRequest) (resp *ports.AuthMethodOtpGetResponse, err error) {
	user, isDeleted, err := s.rel.GetUserByUsername(ctx, req.Username, req.UserType)
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, utils.UserDeletedResponse.Clone()
	}
	if user == nil && req.OtpType != ports.SignupOtpType {
		return nil, utils.UsernameNotFoundResponse.Clone().
			WithReason("username", req.Username)
	}
	if user != nil && req.OtpType == ports.SignupOtpType {
		return nil, utils.UsernameAlreadyExistsResponse.Clone().
			WithReason("username", req.Username)
	}

	var identity string
	if req.OtpType == ports.SignupOtpType {
		if req.Email != "" {
			identity = req.Email
			if err = s.rel.CheckUserEmailLimit(ctx, identity); err != nil {
				return nil, err
			}
		} else {
			identity = req.PhoneNumber
			if err = s.rel.CheckUserPhoneLimit(ctx, identity); err != nil {
				return nil, err
			}
		}
	} else {
		email := user.GetEmail()
		if email != "" {
			identity = email
		} else {
			identity = user.GetPhoneNumber()
		}
	}
	if req.SendToEmail && identity == "" {
		return nil, utils.UserHasNotSetEmailResponse.Clone()
	}
	if !req.SendToEmail && identity == "" {
		return nil, utils.UserHasNotSetPhoneNumberResponse.Clone()
	}

	if req.SendToEmail {
		resp = &ports.AuthMethodOtpGetResponse{MaskedEmail: utils.MaskEmail(identity)}
	} else {
		resp = &ports.AuthMethodOtpGetResponse{MaskedPhone: utils.MaskPhone(identity)}
	}
	prevToken, err := s.cache.GetOtpToken(ctx, identity, req.OtpType, req.UserType)
	if err != nil {
		return nil, err
	}
	if prevToken != "" {
		return nil, utils.AuthMethodOtpGetTooEarlyResponse.Clone()
	}
	token := generateRandomOTP()
	err = s.cache.SetOtpToken(ctx, identity, token, req.OtpType, 2*time.Minute, req.UserType)
	if err != nil {
		return nil, err
	}
	if !req.SendToEmail {
		message := fmt.Sprintf(otpMessages[req.OtpType], token)
		_, err = s.telecom.NoReplySend(ctx, []string{req.PhoneNumber}, message)
		return nil, err
	}
	subject := otpSubjects[req.OtpType]
	message, err := otpMails[req.OtpType](emailStruct{
		Domain:       s.domain,
		Version:      s.version,
		Token:        token,
		SupportEmail: s.supportEmail,
	})
	if err != nil {
		return nil, err
	}
	_, err = s.mailcom.NoReplySend(ctx, []string{user.GetEmail()}, subject, message)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Auth) SendKey(ctx context.Context, req *ports.AuthSignupKeyGetRequest) (err error) {
	key := generateRandomKey()
	if req.Email != "" {
		if err = s.rel.CheckUserEmailLimit(ctx, req.Email); err != nil {
			return err
		}
		if err = s.cache.SetOtpKey(ctx, req.Email, key, 48*time.Hour, req.UserType); err != nil {
			return err
		}
		subject := otpSubjects[ports.AdminSignupKeyOtpType]
		var message string
		message, err = otpMails[ports.AdminSignupKeyOtpType](emailStruct{
			Domain:       s.domain,
			Version:      s.version,
			Token:        key,
			SupportEmail: s.supportEmail,
		})
		if err != nil {
			return err
		}
		_, err = s.mailcom.NoReplySend(ctx, []string{req.Email}, subject, message)

	} else {
		if err = s.rel.CheckUserPhoneLimit(ctx, req.PhoneNumber); err != nil {
			return err
		}
		if err = s.cache.SetOtpKey(ctx, req.PhoneNumber, key, 48*time.Hour, req.UserType); err != nil {
			return err
		}
		message := otpMessages[ports.AdminSignupKeyOtpType]
		_, err = s.telecom.NoReplySend(ctx, []string{req.PhoneNumber}, message)
	}
	return
}

func (s *Auth) GenerateToken(ctx context.Context, user *ports.Login) (err error) {
	defer func() { err = utils.FuncPipe(authCaller+".GenerateToken", err) }()
	now := time.Now()
	user.Jwt = &ports.Jwt{}
	user.Jwt.AccessExpires = int(s.jwtAccessExp / time.Minute)
	user.Jwt.RefreshExpires = int(s.jwtRefreshExp / time.Minute)

	refreshClaims := jwt.MapClaims{
		"id":       user.User.Id,
		"username": user.User.Username,
		"name":     user.User.Name,
		"userType": user.User.UserType,
		"type":     ports.RefreshJwtType,
		"exp":      now.Add(s.jwtRefreshExp).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.jwtSK)
	if err != nil {
		return err
	}
	user.Jwt.RefreshToken = token

	accessClaims := jwt.MapClaims{
		"id":            user.User.Id,
		"username":      user.User.Username,
		"name":          user.User.Name,
		"userType":      user.User.UserType,
		"type":          ports.AccessJwtType,
		"exp":           now.Add(s.jwtAccessExp).Unix(),
		"refresh_token": token,
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.jwtSK)
	if err != nil {
		return err
	}
	user.Jwt.AccessToken = token
	return nil
}

func (s *Auth) ParseToken(token string, jwtType ports.JwtType) (login *ports.Login, err error) {
	defer func() { err = utils.FuncPipe(authCaller+".ParseToken", err) }()
	var claims jwt.MapClaims
	_, err = jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		return s.jwtSK, nil
	})
	if err != nil {
		return nil, utils.JwtUnauthorizedResponse.Clone()
	}
	if claims["type"] != string(jwtType) {
		return nil, utils.JwtUnauthorizedResponse.Clone().
			WithReason("type", claims["type"]).
			WithReason("wanted_type", jwtType)
	}
	id, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		return nil, err
	}
	login = &ports.Login{
		User: &ports.User{
			Id:       id,
			Username: claims["username"].(string),
			Name:     claims["name"].(string),
			UserType: ports.UserType(claims["userType"].(string)),
		},
		Jwt: &ports.Jwt{},
	}
	if jwtType == ports.AccessJwtType {
		login.Jwt.AccessToken = token
		login.Jwt.RefreshToken = claims["refresh_token"].(string)
	} else {
		login.Jwt.RefreshToken = token
	}
	return login, nil
}
