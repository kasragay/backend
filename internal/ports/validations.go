package ports

import (
	"context"
	"encoding/base64"
	"errors"
	"image"
	"regexp"
	"slices"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/utils"
)

var inValidator = validator.New()

func Validate(ctx context.Context, logger *utils.Logger, data any) (err error) {
	defer func() {
		const caller = packageCaller + ".Validate"
		if err != nil {
			var err_ *utils.Error
			if errors.As(err, &err_) {
				err = err_.WithCaller(caller)
			}
		}
	}()
	v := inValidator.Struct(data)
	if v != nil {
		logger.Error(ctx, v, "validation error")
		result, err := utils.StructToMap(data)
		if err != nil {
			return utils.BadRequestResponse.Clone()
		}
		err_ := utils.BadRequestResponse.Clone()
		for key, value := range result {
			if key == "avatar" {
				continue
			}
			err_ = err_.WithReason(key, value)
		}
		return err_
	}
	return nil
}

func init() {
	_ = inValidator.RegisterValidation("nameValidator", nameValidator)
	_ = inValidator.RegisterValidation("usernameValidator", usernameValidator)
	_ = inValidator.RegisterValidation("phoneValidator", phoneValidator)
	_ = inValidator.RegisterValidation("userTypeValidator", userTypeValidator)
	_ = inValidator.RegisterValidation("nonClientUserTypeValidator", nonClientUserTypeValidator)
	_ = inValidator.RegisterValidation("nonAdminUserTypeValidator", nonAdminUserTypeValidator)
	_ = inValidator.RegisterValidation("passwordValidator", passwordValidator)
	_ = inValidator.RegisterValidation("otpTokenValidator", otpTokenValidator)
	_ = inValidator.RegisterValidation("otpKeyValidator", otpKeyValidator)
	_ = inValidator.RegisterValidation("otpTypeValidator", otpTypeValidator)
	_ = inValidator.RegisterValidation("avatarValidator", avatarValidator)
	_ = inValidator.RegisterValidation("avatarObjectNameValidator", avatarObjectNameValidator)
	_ = inValidator.RegisterValidation("emailValidator", emailValidator)

	inValidator.RegisterStructValidation(authMethodOtpGetRequestValidator, AuthMethodOtpGetRequest{})
	inValidator.RegisterStructValidation(authSignupKeyGetRequestValidator, AuthSignupKeyGetRequest{})
	inValidator.RegisterStructValidation(authSignupPostRequestValidator, AuthSignupPostRequest{})
	inValidator.RegisterStructValidation(authSigninPostRequestValidator, AuthSigninPostRequest{})
	inValidator.RegisterStructValidation(authResetPasswordPostRequestValidator, AuthResetPasswordPostRequest{})

}

func nameValidator(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) == 0 {
		return true
	}
	length := len(fl.Field().String())
	if length < 2 || length > 250 {
		return false
	}
	pattern := `^[\x{0600}-\x{06FF}\x{FB50}-\x{FDFF}a-zA-Z\s-]{2,250}$`
	matched, err := regexp.MatchString(pattern, fl.Field().String())
	if err != nil {
		return false
	}
	return matched
}

func usernameValidator(fl validator.FieldLevel) bool {
	length := len(fl.Field().String())
	if length == 0 {
		return true
	}
	if length < 3 || length > 64 {
		return false
	}
	pattern := `^(?!.*\.\.)[a-zA-Z0-9\-_]+(\.[a-zA-Z0-9\-_]+)*$`
	rp := regexp2.MustCompile(pattern, regexp2.None)
	matched, err := rp.MatchString(fl.Field().String())
	if err != nil {
		return false
	}
	return matched
}

func PhoneValidator(phone string) bool {
	if len(phone) != 11 {
		return false
	}
	pattern := `^0\d{10}$`
	matched, err := regexp.MatchString(pattern, phone)
	if err != nil {
		return false
	}
	return matched
}

func phoneValidator(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) == 0 {
		return true
	}
	return PhoneValidator(fl.Field().String())
}

func UserTypeValidator(val string) bool {
	userType := UserType(val)
	return slices.Contains(AllUserTypes, userType)
}

func userTypeValidator(fl validator.FieldLevel) bool {
	return UserTypeValidator(fl.Field().String())
}

func nonClientUserTypeValidator(fl validator.FieldLevel) bool {
	userType := UserType(fl.Field().String())
	if userType == ClientUserType {
		return false
	}
	return slices.Contains(AllUserTypes, userType)
}

func nonAdminUserTypeValidator(fl validator.FieldLevel) bool {
	userType := UserType(fl.Field().String())
	if userType == AdminUserType {
		return false
	}
	return slices.Contains(AllUserTypes, userType)
}

func passwordValidator(fl validator.FieldLevel) bool {
	length := len(fl.Field().String())
	if length == 0 {
		return true
	}
	if length < 8 || length > 100 {
		return false
	}
	// Original regex with lookaheads
	pattern := `^(?=.*[A-Z])(?=.*[0-9])(?=.*[\W_]).{8,100}$`
	re := regexp2.MustCompile(pattern, regexp2.None)
	isMatch, err := re.MatchString(fl.Field().String())
	if err != nil {
		return false
	}
	return isMatch
}

func OtpTokenValidator(value string) bool {
	if len(value) != 5 {
		return false
	}
	pattern := `^\d{5}$`
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false
	}
	return matched
}

func otpTokenValidator(fl validator.FieldLevel) bool {
	return OtpTokenValidator(fl.Field().String())
}

func otpKeyValidator(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) == 0 {
		return true
	}
	if len(fl.Field().String()) != 8 {
		return false
	}
	pattern := `^[a-zA-Z0-9!@#$%^&*()-_=+]{8}$`
	matched, err := regexp.MatchString(pattern, fl.Field().String())
	if err != nil {
		return false
	}
	return matched
}

func otpTypeValidator(fl validator.FieldLevel) bool {
	method := OtpType(fl.Field().String())
	return slices.Contains(OtpTypes, method)
}

func avatarValidator(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	avatar := fl.Field().String()
	_, err := AvatarValidator(avatar)
	return err == nil
}

func AvatarValidator(avatar string) (img *image.Image, err error) {
	if avatar == "" {
		return nil, nil
	}
	dataURIRegex := regexp.MustCompile(`^data:image/(png|jpg|jpeg);base64,(.+)$`)
	matches := dataURIRegex.FindStringSubmatch(avatar)
	if len(matches) != 3 {
		return nil, utils.BadRequestResponse.Clone()
	}
	base64Data := matches[2]
	imgData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, utils.BadRequestResponse.Clone()
	}
	return utils.ImageReader(imgData, "png", [2]int{250, 250}, []string{"png", "jpg", "jpeg"})
}

func avatarObjectNameValidator(fl validator.FieldLevel) bool {
	avatarObjectName := fl.Field().String()
	if !strings.HasSuffix(avatarObjectName, ".png") {
		return false
	}
	avatarObjectId := strings.TrimSuffix(avatarObjectName, ".png")
	_, err := uuid.Parse(avatarObjectId)
	return err == nil
}

func EmailValidator(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}
	return matched
}

func emailValidator(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) == 0 {
		return true
	}
	return EmailValidator(fl.Field().String())
}

func authMethodOtpGetRequestValidator(sl validator.StructLevel) {
	req := sl.Current().Interface().(AuthMethodOtpGetRequest)

	if req.SendToEmail == req.SendToPhone {
		sl.ReportError(req.SendToEmail, "SendToEmail,SendToPhone", "", "exactlyOneSendMethod", "")
	}
	if (req.Email != "" && req.PhoneNumber != "") || (req.Email == "" && req.PhoneNumber == "") {
		sl.ReportError(req.Email, "Email,PhoneNumber", "", "exactlyOneSendMethod", "")
	}
	if !req.SendToEmail && req.Email != "" {
		sl.ReportError(req.Email, "Email", "", "emailMustBeEmpty", "")
	}
	if req.SendToEmail && req.Email == "" && req.OtpType == SignupOtpType {
		sl.ReportError(req.Email, "Email", "", "emailIsRequired", "")
	}
	if !req.SendToPhone && req.PhoneNumber != "" {
		sl.ReportError(req.PhoneNumber, "PhoneNumber", "", "phoneNumberMustBeEmpty", "")
	}
	if req.SendToPhone && req.PhoneNumber == "" && req.OtpType == SignupOtpType {
		sl.ReportError(req.PhoneNumber, "PhoneNumber", "", "phoneNumberIsRequired", "")
	}

}

func authSignupKeyGetRequestValidator(sl validator.StructLevel) {
	req := sl.Current().Interface().(AuthSignupKeyGetRequest)
	if (req.Email != "" && req.PhoneNumber != "") || (req.Email == "" && req.PhoneNumber == "") {
		sl.ReportError(req.Email, "Email,PhoneNumber", "", "exactlyOneSendMethod", "")
	}
}

func authSignupPostRequestValidator(sl validator.StructLevel) {
	req := sl.Current().Interface().(AuthSignupPostRequest)
	if (req.Email != "" && req.PhoneNumber != "") || (req.Email == "" && req.PhoneNumber == "") {
		sl.ReportError(req.Email, "Email,PhoneNumber", "", "exactlyOneSendMethod", "")
	}
}
func authSigninPostRequestValidator(sl validator.StructLevel) {
	req := sl.Current().Interface().(AuthSigninPostRequest)
	if req.SentToEmail == req.SentToPhone {
		sl.ReportError(req.SentToEmail, "SentToEmail,SentToPhone", "", "exactlyOneSendMethod", "")
	}
}

func authResetPasswordPostRequestValidator(sl validator.StructLevel) {
	req := sl.Current().Interface().(AuthResetPasswordPostRequest)
	if req.SentToEmail == req.SentToPhone {
		sl.ReportError(req.SentToEmail, "SentToEmail,SentToPhone", "", "exactlyOneSendMethod", "")
	}
}
