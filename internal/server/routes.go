package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kasragay/backend/internal/ports"
)

func (s *AbstractServer) registerBasicRoutes() {

	s.App().Use(s.LoggerMiddleware())
	s.App().Use(cors.New(cors.Config{
		AllowOrigins:     fmt.Sprintf("%s,%s", s.BackUrl(), s.FrontUrl()),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Accept,Authorization,Content-Type",
		AllowCredentials: true,
		MaxAge:           300,
	}))
	switch s.serviceName {
	case ports.GatewayServiceName:
		s.App().Use(s.allowedHostsMiddleware("api." + s.Domain()))
	default:
		s.App().Use(s.allowedHostsMiddleware(string(s.serviceName)))
		s.App().Use(s.authMiddleware_)
	}
	s.verRouter = s.App().Group("/" + string(s.Version()))
}
