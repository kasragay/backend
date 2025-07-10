package gateway

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
)

func (s *GatewayServer) AuthUrlMiddleware(jwtType ports.JwtType, onlyAdmin bool) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			const caller = gatewayServerCaller + ".AuthUrlMiddleware"
			if err != nil {
				var uErr *utils.Error
				if errors.As(err, &uErr) {
					err = uErr.WithCaller(caller)
				} else {
					err = utils.NewInternalError(err).WithCaller(caller)
				}
			}
		}()
		return s.authMiddleware(c.Params("jwt"), jwtType, onlyAdmin)(c)
	}
}

func (s *GatewayServer) AuthBearerMiddleware(jwtType ports.JwtType, onlyAdmin bool) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			const caller = gatewayServerCaller + ".AuthBearerMiddleware"
			if err != nil {
				var uErr *utils.Error
				if errors.As(err, &uErr) {
					err = uErr.WithCaller(caller)
				} else {
					err = utils.NewInternalError(err).WithCaller(caller)
				}
			}
		}()
		token := c.Get(fiber.HeaderAuthorization)
		if token == "" {
			return utils.JwtUnauthorizedResponse.Clone()
		}
		token = strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
		return s.authMiddleware(token, jwtType, onlyAdmin)(c)
	}
}

func (s *GatewayServer) authMiddleware(token string, jwtType ports.JwtType, onlyAdmin bool) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			const caller = gatewayServerCaller + ".authMiddleware"
			if err != nil {
				var uErr *utils.Error
				if errors.As(err, &uErr) {
					err = uErr.WithCaller(caller)
				} else {
					err = utils.NewInternalError(err).WithCaller(caller)
				}
			}
		}()
		if token == "" {
			return utils.JwtUnauthorizedResponse.Clone()
		}
		login, err := s.auth.ParseToken(token, jwtType)
		if err != nil {
			return err
		}
		if onlyAdmin && login.User.UserType != ports.AdminUserType {
			return utils.JwtUnauthorizedResponse.Clone()
		}
		if isIn, err := s.Cache().IsJwtInBlacklist(c.Context(), token); err != nil {
			return err
		} else if isIn {
			return utils.JwtUnauthorizedResponse.Clone()
		}
		if chResp, isDeleted, err := s.auth.CheckById(c.Context(), login.User.Id, login.User.UserType); err != nil {
			return err
		} else if !chResp || isDeleted {
			return utils.JwtUnauthorizedResponse.Clone()
		}
		c.Request().Header.Set("X-User-ID", login.User.Id.String())
		c.Request().Header.Set("X-User-Username", login.User.Username)
		c.Request().Header.Set("X-User-UserType", string(login.User.UserType))
		c.Locals("login", login)
		c.Locals("id", login.User.Id)
		c.Locals("username", login.User.Username)
		c.Locals("userType", login.User.UserType)
		return c.Next()
	}
}
