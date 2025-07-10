package gateway

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
	"github.com/minio/minio-go/v7"
	"github.com/valyala/fasthttp"
)

func (s *GatewayServer) RegisterRoutes() {
	s.App().Use(favicon.New(favicon.Config{
		File: "./assets/favicon.ico",
		URL:  "/favicon.ico",
	}))
	switch s.Version() {
	case ports.V1Version:
		s.RegisterV1()
	}
}

func (s *GatewayServer) RegisterV1() {
	s.register(
		"/", s.App(), ports.GET, "/", s.redirectToLatestVersionHandler,
		20, time.Minute, true, false, false,
	)
	s.register(
		string(s.Version()), s.VersionRouter(), ports.GET, "/", s.redirectToDocsHandler,
		20, time.Minute, true, false, false,
	)
	s.register(
		string(s.Version()), s.VersionRouter(), ports.GET, "/health", s.healthGetHandler,
		20, time.Minute, true, false, false,
	)
	s.VersionRouter().Static("/docs", "./docs")
	s.VersionRouter().Static("/assets", "./assets")
	s.register(
		string(s.Version()), s.VersionRouter(), ports.GET, "/docs/*", NewSwaggerUI(
			SwaggerUIConfig{
				URL: fmt.Sprintf("/%s/docs/swagger.yaml", string(s.Version())),
			},
		),
		20, time.Minute, true, false, false,
	)
	s.register(
		string(s.Version()), s.VersionRouter(), ports.GET, "/s3/avatars/:userType/:objectName", s.s3AvatarsGetHandler,
		20, time.Minute, true, false, false,
	)
	auth := s.VersionRouter().Group("/auth")
	s.register(
		"auth", auth, ports.GET, "/check", s.authCheckPostHandler,
		10, time.Minute, true, false, false,
	)
	s.register(
		"auth", auth, ports.GET, "/tmp/signup/key", s.tmpAuthSignupKeyGetHandler,
		3, time.Minute, false, true, true,
		s.AuthBearerMiddleware(ports.AccessJwtType, true),
	)
	s.register(
		"auth", auth, ports.GET, "/signup/key", s.authSignupKeyGetHandler,
		3, time.Minute, false, true, true,
		s.AuthBearerMiddleware(ports.AccessJwtType, true),
	)
	s.register(
		"auth", auth, ports.GET, "/tmp/:otpType/otp", s.tmpAuthMethodOtpGetHandler,
		3, time.Minute, true, false, false,
	)
	s.register(
		"auth", auth, ports.GET, "/:otpType/otp", s.authMethodOtpGetHandler,
		3, time.Minute, true, false, false,
	)
	s.register(
		"auth", auth, ports.POST, "/signup", s.authSignupPostHandler,
		10, time.Minute, true, false, false,
		s.ContentTypeMiddleware(fiber.MIMEApplicationJSON),
	)
	s.register(
		"auth", auth, ports.POST, "/signin", s.authSigninPostHandler,
		10, time.Minute, true, false, false,
		s.ContentTypeMiddleware(fiber.MIMEApplicationJSON),
	)
	s.register(
		"auth", auth, ports.POST, "/signin/password", s.authSigninPasswordPostHandler,
		10, time.Minute, true, false, false,
		s.ContentTypeMiddleware(fiber.MIMEApplicationJSON),
	)
	s.register(
		"auth", auth, ports.POST, "/logout", s.authLogoutPostHandler,
		10, time.Minute, false, true, true,
		s.AuthBearerMiddleware(ports.AccessJwtType, false),
	)
	s.register(
		"auth", auth, ports.POST, "/refresh", s.authRefreshPostHandler,
		10, time.Minute, false, true, true,
		s.AuthBearerMiddleware(ports.RefreshJwtType, false),
	)
	s.register(
		"auth", auth, ports.POST, "/reset-password", s.authResetPasswordPostHandler,
		5, time.Minute, true, false, false,
		s.ContentTypeMiddleware(fiber.MIMEApplicationJSON),
	)
	s.register(
		"auth", auth, ports.POST, "/reset-phone", s.authResetPhonePostHandler,
		5, time.Minute, false, true, true,
		s.AuthBearerMiddleware(ports.AccessJwtType, false),
		s.ContentTypeMiddleware(fiber.MIMEApplicationJSON),
	)
	s.register(
		"auth", auth, ports.POST, "/reset-email", s.authResetEmailPostHandler,
		5, time.Minute, false, true, true,
		s.AuthBearerMiddleware(ports.AccessJwtType, false),
		s.ContentTypeMiddleware(fiber.MIMEApplicationJSON),
	)

	user := s.VersionRouter().Group("/user")
	s.register(
		"user", user, ports.GET, "/", s.userServiceProxyHandler,
		10, time.Minute, false, true, true,
		s.AuthBearerMiddleware(ports.AccessJwtType, false),
	)
	s.register(
		"user", user, ports.PUT, "/", s.userServiceProxyHandler,
		10, time.Minute, false, true, true,
		s.AuthBearerMiddleware(ports.AccessJwtType, false),
	)
	s.register(
		"user", user, ports.DELETE, "/", s.userServiceProxyHandler,
		10, time.Minute, false, true, true,
		s.AuthBearerMiddleware(ports.AccessJwtType, false),
	)

	client := s.VersionRouter().Group("/client")
	_ = client
}

func (s *GatewayServer) register(
	routerName string,
	router fiber.Router,
	method ports.Method,
	path string,
	handler fiber.Handler,
	token uint64,
	duration time.Duration,
	onIP, onID, hasAuthMiddleware bool,
	middlewares ...fiber.Handler,
) {
	if hasAuthMiddleware && len(middlewares) == 0 {
		s.Logger().Fatalf(context.Background(), "auth middleware is not set for: %s", method)
	}
	token = token * 10 // TODO_DEL
	newMiddlewares := make([]fiber.Handler, 0, len(middlewares)+2)
	ratelimiter := s.ratelimiter.Handler(routerName, method, path, token, duration, onIP, onID)
	if hasAuthMiddleware {
		newMiddlewares = append(newMiddlewares, middlewares[0])
		newMiddlewares = append(newMiddlewares, ratelimiter)
		if len(middlewares) > 1 {
			newMiddlewares = append(newMiddlewares, middlewares[1:]...)
		}
		newMiddlewares = append(newMiddlewares, handler)
	} else {
		newMiddlewares = append(newMiddlewares, ratelimiter)
		newMiddlewares = append(newMiddlewares, middlewares...)
		newMiddlewares = append(newMiddlewares, handler)
	}
	switch method {
	case ports.GET:
		router.Get(path, newMiddlewares...)
	case ports.POST:
		router.Post(path, newMiddlewares...)
	case ports.PUT:
		router.Put(path, newMiddlewares...)
	case ports.DELETE:
		router.Delete(path, newMiddlewares...)
	case ports.OPTIONS:
		router.Options(path, newMiddlewares...)
	case ports.PATCH:
		router.Patch(path, newMiddlewares...)
	default:
		s.Logger().Fatalf(context.Background(), "method is not valid for: %s", method)
		return
	}
}

func (s *GatewayServer) healthGetHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".healthGetHandler", err) }()
	return c.SendStatus(fiber.StatusOK)
}

func (s *GatewayServer) redirectToLatestVersionHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".redirectToLatestVersionHandler", err) }()
	return c.Redirect(fmt.Sprintf("/%s", string(ports.LatestVersion)))
}

func (s *GatewayServer) redirectToDocsHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".redirectToDocsHandler", err) }()
	return c.Redirect(fmt.Sprintf("/%s/docs", string(s.Version())))
}

func (s *GatewayServer) s3AvatarsGetHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".s3AvatarsGetHandler", err) }()
	req := ports.S3AvatarsGetRequest{
		UserType:   ports.UserType(c.Params("userType")),
		ObjectName: c.Params("objectName"),
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	id, err := uuid.Parse(strings.TrimSuffix(req.ObjectName, ".png"))
	if err != nil {
		return utils.BadRequestResponse.Clone().
			WithReason("object_name", req.ObjectName)
	}
	avatar, err := s.S3().GetAvatar(c.Context(), id, req.UserType)
	if err != nil {
		return err
	}
	if avatar == nil {
		return utils.NotFoundResponse.Clone()
	}
	defer func() {
		_ = avatar.Close()
	}()
	objInfo, err := avatar.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return utils.NotFoundResponse.Clone()
		}
		return err
	}
	contentType := objInfo.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Set(fiber.HeaderContentType, contentType)
	c.Set(fiber.HeaderContentLength, strconv.FormatInt(objInfo.Size, 10))
	c.Set(fiber.HeaderCacheControl, "private, max-age=0") // 1 year
	buffer := make([]byte, 32*1024)
	for {
		n, err := avatar.Read(buffer)
		if n > 0 {
			if _, writeErr := c.Response().BodyWriter().Write(buffer[:n]); writeErr != nil {
				return writeErr
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return
}

func (s *GatewayServer) authCheckPostHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".authCheckPostHandler", err) }()
	req := ports.AuthCheckPostRequest{
		Username: c.Query("username"),
		UserType: ports.UserType(c.Query("user_type")),
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	resp, err := s.auth.CheckPost(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)

}

func (s *GatewayServer) tmpAuthSignupKeyGetHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".tmpAuthAdminSignupKeyGetHandler", err) }()
	req := ports.AuthSignupKeyGetRequest{
		Email:       c.Query("email"),
		PhoneNumber: c.Query("phone_number"),
		UserType:    ports.UserType(c.Query("user_type")),
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	resp, err := s.auth.TmpSignupKeyGet(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *GatewayServer) authSignupKeyGetHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".authAdminSignupKeyGetHandler", err) }()
	req := ports.AuthSignupKeyGetRequest{
		Email:       c.Query("email"),
		PhoneNumber: c.Query("phone_number"),
		UserType:    ports.UserType(c.Query("user_type")),
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	return s.auth.SignupKeyGet(c.Context(), &req)
}

func (s *GatewayServer) tmpAuthMethodOtpGetHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".tmpAuthMethodOtpGetHandler", err) }()
	stre := c.Query("send_to_email")
	if stre == "" {
		stre = "false"
	}
	parsedStre, err := strconv.ParseBool(stre)
	if err != nil {
		return utils.BadRequestResponse.Clone().WithReason("send_to_email", stre)
	}
	strp := c.Query("send_to_phone")
	if strp == "" {
		strp = "false"
	}
	parsedStrp, err := strconv.ParseBool(strp)
	if err != nil {
		return utils.BadRequestResponse.Clone().WithReason("send_to_phone", strp)
	}
	req := ports.AuthMethodOtpGetRequest{
		Username:    c.Query("username"),
		SendToEmail: parsedStre,
		SendToPhone: parsedStrp,
		OtpType:     ports.OtpType(c.Params("otpType")),
		UserType:    ports.UserType(c.Query("user_type")),
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	resp, err := s.auth.TmpMethodOtpGet(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *GatewayServer) authMethodOtpGetHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".authMethodOtpGetHandler", err) }()
	stre := c.Query("send_to_email")
	if stre == "" {
		stre = "false"
	}
	parsedStre, err := strconv.ParseBool(stre)
	if err != nil {
		return utils.BadRequestResponse.Clone().WithReason("send_to_email", stre)
	}
	strp := c.Query("send_to_phone")
	if strp == "" {
		strp = "false"
	}
	parsedStrp, err := strconv.ParseBool(strp)
	if err != nil {
		return utils.BadRequestResponse.Clone().WithReason("send_to_phone", strp)
	}
	req := ports.AuthMethodOtpGetRequest{
		Username:    c.Query("username"),
		SendToEmail: parsedStre,
		SendToPhone: parsedStrp,
		OtpType:     ports.OtpType(c.Params("otpType")),
		UserType:    ports.UserType(c.Query("user_type")),
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	resp, err := s.auth.MethodOtpGet(c.Context(), &req)
	if err != nil {
		return err
	}
	if resp == nil {
		return nil
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *GatewayServer) authSignupPostHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".authSignupPostHandler", err) }()
	req := ports.AuthSignupPostRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse.Clone()
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	resp, err := s.auth.SignupPost(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *GatewayServer) authSigninPostHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".authSigninPostHandler", err) }()
	req := ports.AuthSigninPostRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse.Clone()
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	resp, err := s.auth.SigninPost(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *GatewayServer) authSigninPasswordPostHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".authSigninPasswordPostHandler", err) }()
	req := ports.AuthSigninPasswordPostRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse.Clone()
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	resp, err := s.auth.SigninPasswordPost(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *GatewayServer) authLogoutPostHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".authLogoutPostHandler", err) }()
	login := c.Locals("login").(*ports.Login)
	return s.auth.LogoutPost(c.Context(), login)

}

func (s *GatewayServer) authRefreshPostHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".authRefreshPostHandler", err) }()
	login := c.Locals("login").(*ports.Login)
	resp, err := s.auth.RefreshPost(c.Context(), login)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *GatewayServer) authResetPasswordPostHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".authResetPasswordPostHandler", err) }()
	req := ports.AuthResetPasswordPostRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse.Clone()
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}

	resp, err := s.auth.ResetPasswordPost(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *GatewayServer) authResetPhonePostHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".authResetPhonePostHandler", err) }()
	req := ports.AuthResetPhonePostRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse.Clone()
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	if err := s.CanPass(c, req.Id); err != nil {
		return err
	}
	return s.auth.ResetPhonePost(c.Context(), &req)
}

func (s *GatewayServer) authResetEmailPostHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".authResetEmailPostHandler", err) }()
	req := ports.AuthResetEmailPostRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse.Clone()
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	if err := s.CanPass(c, req.Id); err != nil {
		return err
	}
	return s.auth.ResetEmailPost(c.Context(), &req)
}

func (s *GatewayServer) userServiceProxyHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(gatewayServerCaller+".userServiceProxyHandler", err) }()
	client := &fasthttp.Client{
		MaxConnWaitTimeout: time.Second,
		ReadTimeout:        2 * time.Second,
		WriteTimeout:       2 * time.Second,
	}
	return proxy.Forward("http://user:8082"+c.OriginalURL(), client)(c)
}
