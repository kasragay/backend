package post

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
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
	// TODO: change handlers to handle Id param instead of reading id from body
	post.Get("/:id", s.postGetHandler)
	post.Post("/", s.postPostHandler)
	post.Put("/:id", s.ContentTypeMiddleware(fiber.MIMEApplicationJSON), s.postPutHandler)
	post.Delete("/:id", s.ContentTypeMiddleware(fiber.MIMEApplicationJSON), s.postDeleteHandler)

	_ = s.VersionRouter().Group("/client")

	//tag := s.VersionRouter().Group("/tag")
	// TODO: implement tag endpoints
}

func (s *PostServer) postHealthGetHandler(c *fiber.Ctx) (err error) {
	return c.SendStatus(fiber.StatusOK)
}

// TODO: implement handler methods

func (s *PostServer) postGetHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(postServerCaller+".postGetHandler", err) }()
	parsedId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse.Clone().
			WithReason("id", c.Query("id"))
	}
	if err := s.CanPass(c, parsedId); err != nil {
		return err
	}

	req := ports.PostGetRequest{
		Id: parsedId,
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	resp, err := s.post.PostGet(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)

}
func (s *PostServer) postPostHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(postServerCaller+".postPostHandler", err) }()
	req := ports.PostPostRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse.Clone()
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	resp, err := s.post.PostPost(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}
func (s *PostServer) postPutHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(postServerCaller+".postPutHandler", err) }()
	req := ports.PostPutRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse.Clone()
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
	}
	if err := s.CanPass(c, req.Id); err != nil {
		return err
	}
	err = s.post.PostPut(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusAccepted)
}
func (s *PostServer) postDeleteHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(postServerCaller+".postDeleteHandler", err) }()
	parsedId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequestResponse.Clone().
			WithReason("id", c.Query("id"))
	}
	if err := s.CanPass(c, parsedId); err != nil {
		return err
	}

	err = s.post.PostDelete(c.Context(), parsedId)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}
