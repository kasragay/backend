package server

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
	"go.uber.org/zap"
)

func (s *AbstractServer) LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		start := time.Now()
		regex := regexp.MustCompile(fmt.Sprintf(`^/%s/[a-zA-Z0-9_-]+/health$`, string(s.Version())))
		if regex.MatchString(c.Path()) {
			return nil
		}
		if err = c.Next(); err != nil {
			s.logger.ErrorFields(
				c.Context(),
				err,
				fmt.Sprintf("Request took %s", time.Since(start)),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
			return err
		}
		s.logger.InfoFields(
			c.Context(),
			fmt.Sprintf("Request took %s", time.Since(start)),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)
		return nil
	}
}

func (s *AbstractServer) allowedHostsMiddleware(allowedHosts ...string) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		host := c.Get(fiber.HeaderHost)
		if idx := strings.Index(host, ":"); idx != -1 {
			host = host[:idx]
		}
		for _, allowedHost := range allowedHosts {
			if strings.EqualFold(host, allowedHost) {
				return c.Next()
			}
		}
		return utils.NewError(fiber.StatusForbidden, "Host not allowed").
			WithReason("host", c.Get(fiber.HeaderHost)).
			WithCaller(serverCaller + ".allowedHostsMiddleware")
	}
}
func (s *AbstractServer) authMiddleware_(c *fiber.Ctx) (err error) {
	defer func() {
		if err != nil {
			var uErr *utils.Error
			if errors.As(err, &uErr) {
				err = uErr.WithCaller(serverCaller + ".authMiddleware_")
			}
		}
	}()
	re := regexp.MustCompile((fmt.Sprintf(`^/%s/([a-zA-Z0-9_-]+)/health$`, string(s.Version()))))
	if re.MatchString(c.OriginalURL()) {
		return c.Next()
	}
	sourceId, err := uuid.Parse(c.Get("X-User-ID"))
	if err != nil {
		return utils.BadRequestResponse.Clone().
			WithReason("X-User-ID", c.Get("X-User-ID"))
	}
	sourcePhone := c.Get("X-User-Phone")
	if sourcePhone == "" {
		return utils.BadRequestResponse.Clone().
			WithReason("X-User-Phone", c.Get("X-User-Phone"))
	}
	sourceUserType := ports.UserType(c.Get("X-User-UserType"))
	if sourceUserType == "" {
		return utils.BadRequestResponse.Clone().
			WithReason("X-User-UserType", c.Get("X-User-UserType"))
	}
	c.Locals("id", sourceId)
	c.Locals("phone", sourcePhone)
	c.Locals("userType", sourceUserType)
	return c.Next()
}

func (s *AbstractServer) ContentTypeMiddleware(acceptedContentTypes ...string) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		if !slices.Contains(acceptedContentTypes, c.Get("Content-Type")) {
			return utils.UnsupportedMediaTypeResponse.Clone().
				WithReason("Content-Type", c.Get("Content-Type"))
		}
		return c.Next()
	}
}
