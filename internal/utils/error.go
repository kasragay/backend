package utils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"maps"

	"github.com/gofiber/fiber/v2"
)

var Debug bool

func init() {
	debug := os.Getenv("DEBUG")
	if debug == "" {
		debug = "false"
	}
	var err error
	Debug, err = strconv.ParseBool(debug)
	if err != nil {
		panic(err)
	}
}

type AppCode int

const (
	NotImplementedAppCode AppCode = 998 + iota
	TooManyRequestsAppCode
	InternalServerAppCode
	NotFoundAppCode
	InvalidHostAppCode
	UnsupportedMediaTypeAppCode
	BadRequestAppCode
	UsernameAlreadyExistsAppCode
	UserNotFoundAppCode
	UsernameNotFoundAppCode
	TokenIncorrectAppCode
	PasswordIncorrectAppCode
	JwtUnauthorizedAppCode
	AuthMethodOtpGetTooEarlyAppCode
	AuthSignupKeyIncorrectAppCode
	BadAvatarAppCode
	UsernameIsTakenAppCode
	UserDeletedAppCode
	UserHasNotSetEmailAppCode
	UserHasNotSetPhoneNumberAppCode
	EmailIsTakenByMultipleAppCode
	PhoneNumberIsTakenByMultipleAppCode
	PostDeletedAppCode
)

var (
	NotImplementedResponse               = NewError(http.StatusNotImplemented, "not implemented").WithAppCode(NotImplementedAppCode)
	TooManyRequestsResponse              = NewError(http.StatusTooManyRequests, "too many requests").WithAppCode(TooManyRequestsAppCode)
	InternalServerResponse               = NewError(http.StatusInternalServerError, "internal server error").WithAppCode(InternalServerAppCode)
	NotFoundResponse                     = NewError(http.StatusNotFound, "not found").WithAppCode(NotFoundAppCode)
	InvalidHostResponse                  = NewError(http.StatusBadRequest, "invalid host").WithAppCode(InvalidHostAppCode)
	UnsupportedMediaTypeResponse         = NewError(http.StatusUnsupportedMediaType, "unsupported media type").WithAppCode(UnsupportedMediaTypeAppCode)
	BadRequestResponse                   = NewError(http.StatusBadRequest, "bad request").WithAppCode(BadRequestAppCode)
	UsernameAlreadyExistsResponse        = NewError(http.StatusConflict, "username already exists").WithAppCode(UsernameAlreadyExistsAppCode)
	UserNotFoundResponse                 = NewError(http.StatusNotFound, "user not found").WithAppCode(UserNotFoundAppCode)
	UsernameNotFoundResponse             = NewError(http.StatusNotFound, "username not found").WithAppCode(UsernameNotFoundAppCode)
	TokenIncorrectResponse               = NewError(http.StatusUnauthorized, "token incorrect").WithAppCode(TokenIncorrectAppCode)
	PasswordIncorrectResponse            = NewError(http.StatusUnauthorized, "password incorrect").WithAppCode(PasswordIncorrectAppCode)
	JwtUnauthorizedResponse              = NewError(http.StatusUnauthorized, "jwt unauthorized").WithAppCode(JwtUnauthorizedAppCode)
	AuthMethodOtpGetTooEarlyResponse     = NewError(http.StatusTooManyRequests, "otp get is called too early").WithAppCode(AuthMethodOtpGetTooEarlyAppCode) // do not change this errpr's body because of ratelimiter
	AuthSignupKeyIncorrectResponse       = NewError(http.StatusUnauthorized, "signup key incorrect").WithAppCode(AuthSignupKeyIncorrectAppCode)
	BadAvatarResponse                    = NewError(http.StatusBadRequest, "bad avatar image").WithAppCode(BadAvatarAppCode)
	UsernameIsTakenResponse              = NewError(http.StatusConflict, "username is taken").WithAppCode(UsernameIsTakenAppCode)
	UserDeletedResponse                  = NewError(http.StatusGone, "user has been deleted").WithAppCode(UserDeletedAppCode)
	UserHasNotSetEmailResponse           = NewError(http.StatusBadRequest, "user has not set email").WithAppCode(UserHasNotSetEmailAppCode)
	UserHasNotSetPhoneNumberResponse     = NewError(http.StatusBadRequest, "user has not set phone number").WithAppCode(UserHasNotSetPhoneNumberAppCode)
	EmailIsTakenByMultipleResponse       = NewError(http.StatusConflict, "email is taken by multiple users").WithAppCode(EmailIsTakenByMultipleAppCode)
	PhoneNumberIsTakenByMultipleResponse = NewError(http.StatusConflict, "phone number is taken by multiple users").WithAppCode(PhoneNumberIsTakenByMultipleAppCode)
	PostDeletedResponse                  = NewError(http.StatusGone, "post has been deleted").WithAppCode(PostDeletedAppCode)
)

type Error struct {
	status     string
	code       int
	appCode    AppCode
	message    string
	reasons    map[string]any
	callers    []string
	isInternal bool
}

func (e Error) IsInternal() bool {
	return e.isInternal
}

func (e Error) GetAppCode() AppCode {
	return e.appCode
}

func (e Error) GetCode() int {
	return e.code
}

func (e Error) GetMessage() string {
	return e.message
}

func (e Error) GetReasons() map[string]any {
	return e.reasons
}

func (e Error) GetCallers() []string {
	return e.callers
}

// NewError returns an Error with a status of "error", the provided code and
// message, and empty reasons and callers. The message is trimmed of trailing
// punctuation characters.
func NewError(code int, message string) *Error {
	message = strings.TrimRight(message, `!?.`)
	return &Error{
		status:  "error",
		code:    code,
		message: message,
		callers: make([]string, 0),
		reasons: make(map[string]any),
	}
}

func FuncPipe(caller string, err error) (resErr error) {
	if err != nil {
		var uErr *Error
		if errors.As(err, &uErr) {
			resErr = uErr.WithCaller(caller)
		} else {
			resErr = NewInternalError(err).WithCaller(caller)
		}
	}
	return
}

func NewInternalError(err error, message ...string) *Error {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	if msg == "" {
		msg = err.Error()
	} else {
		msg += ": " + err.Error()
	}
	code := int(InternalServerAppCode)
	msg = strings.TrimRight(msg, `!?.`)
	return &Error{
		status:     "error",
		code:       code,
		message:    msg,
		callers:    make([]string, 0),
		reasons:    make(map[string]any),
		isInternal: true,
	}
}

func (e *Error) Clone() *Error {
	reasons := make(map[string]any)
	maps.Copy(reasons, e.reasons)
	callers := make([]string, len(e.callers))
	copy(callers, e.callers)
	return &Error{
		status:     e.status,
		code:       e.code,
		appCode:    e.appCode,
		message:    e.message,
		reasons:    reasons,
		callers:    callers,
		isInternal: e.isInternal,
	}
}

// WithReason adds a reason to the Error with the given key and value.
// It returns the same Error for chaining.
func (e *Error) WithReason(key string, value any) *Error {
	e.reasons[key] = value
	return e
}

// WithAppCode adds the provided AppCode to the Error. It returns the same Error
// for chaining.
func (e *Error) WithAppCode(code AppCode) *Error {
	e.appCode = code
	return e
}

// WithCaller adds the provided caller to the Error's list of callers. It
// returns the same Error for chaining.
func (e *Error) WithCaller(caller string) *Error {
	e.callers = append([]string{caller}, e.callers...)
	return e
}

func ErrorHandlerFunc(logger *Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var fe *fiber.Error
		if errors.As(err, &fe) {
			switch fe.Code {
			case fiber.StatusNotFound:
				return c.Status(fe.Code).JSON(NotFoundResponse.getErrorResp())
			case fiber.StatusInternalServerError:
				return c.Status(fe.Code).JSON(InternalServerResponse.getErrorResp())
			}
		}
		var e *Error
		if !errors.As(err, &e) {
			logger.Errorf(c.Context(), err, "##unknown##")
			return c.Status(fiber.StatusInternalServerError).JSON(InternalServerResponse.getErrorResp())
		}
		conErr := err.(*Error)
		if conErr.isInternal {
			logger.Errorf(c.Context(), err, "##internal##")
			return c.Status(fiber.StatusInternalServerError).JSON(InternalServerResponse.getErrorResp())
		}
		return c.Status(conErr.code).JSON(conErr.getErrorResp())
	}
}

// Error returns the error message as a string, making the Error struct
// compatible with the error interface.
func (e Error) Error() string {
	return e.message
}

// getErrorResp returns the Error struct as an ErrorResp,
// which is the struct that is marshaled to JSON and sent as the response body
// for errors.
func (e Error) getErrorResp() ErrorResp {
	if Debug {
		return ErrorResp{
			Status:  e.status,
			Code:    int(e.appCode),
			Message: e.message,
			Reasons: e.reasons,
			Callers: e.callers,
		}
	}
	return ErrorResp{
		Status:  e.status,
		Code:    int(e.appCode),
		Message: e.message,
		Reasons: e.reasons,
	}
}

type ErrorResp struct {
	Status  string         `json:"status"`
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Reasons map[string]any `json:"reasons,omitempty"`
	Callers []string       `json:"callers,omitempty"`
}

func (e *ErrorResp) ToError(code int) *Error {
	return &Error{
		status:  e.Status,
		code:    code,
		appCode: AppCode(e.Code),
		message: e.Message,
		reasons: e.Reasons,
		callers: e.Callers,
	}
}

func HandleUtilsError(logger *Logger, ctx context.Context, w http.ResponseWriter, err *Error) bool {
	conErr := err.Clone().WithCaller("utils.HandleUtilsError")
	if conErr.isInternal {
		logger.Errorf(ctx, err, "##internal##")
		conErr = InternalServerResponse.Clone()
	}
	w.WriteHeader(conErr.GetCode())
	if encodeErr := json.NewEncoder(w).Encode(conErr.getErrorResp()); encodeErr != nil {
		logger.Errorf(ctx, encodeErr, "failed to encode error response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(InternalServerResponse.getErrorResp())
	}
	return true
}

// HandleUnknownError writes an unknown error to an http.ResponseWriter as a utils.Error with InternalServerError status.
// It logs the error with a reason and returns true to indicate the error was handled.
func HandleUnknownError(logger *Logger, ctx context.Context, w http.ResponseWriter, err error) bool {
	conErr := InternalServerResponse.Clone().WithCaller("utils.HandleUnknownError")
	logger.Error(ctx, err, "##unknown##")
	w.WriteHeader(conErr.GetCode())
	if encodeErr := json.NewEncoder(w).Encode(conErr.getErrorResp()); encodeErr != nil {
		logger.Errorf(ctx, encodeErr, "failed to encode error response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(InternalServerResponse.getErrorResp())
	}
	return true
}
