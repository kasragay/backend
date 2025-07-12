package post

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kasragay/backend/internal/ports"
)

func (s *PostServer) RegisterRoutes() {
	switch s.Version() {
	case ports.V1Version:
		s.RegisterV1()
	}
}

func (s *PostServer) RegisterV1() {
	post := s.VersionRouter().Group("/post")
	post.Get("/health", s.postHealthGetHandler)

	post.Get("/", s.postGetHandler)
	post.Put("/", s.ContentTypeMiddleware(fiber.MIMEApplicationJSON), s.postPutHandler)
	post.Delete("/", s.ContentTypeMiddleware(fiber.MIMEApplicationJSON), s.postDeleteHandler)

	client := s.VersionRouter().Group("/client")
	_ = client
}

func (s *PostServer) postHealthGetHandler(c *fiber.Ctx) (err error) {
	return c.SendStatus(fiber.StatusOK)
}

// TODO: implement handler methods

func (s *PostServer) postGetHandler(c *fiber.Ctx) (err error)
func (s *PostServer) postPutHandler(c *fiber.Ctx) (err error)
func (s *PostServer) postPostHandler(c *fiber.Ctx) (err error)
func (s *PostServer) postDeleteHandler(c *fiber.Ctx) (err error)
