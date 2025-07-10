package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
)

func (s *UserServer) RegisterRoutes() {
	switch s.Version() {
	case ports.V1Version:
		s.RegisterV1()
	}
}

func (s *UserServer) RegisterV1() {
	user := s.VersionRouter().Group("/user")
	user.Get("/health", s.userHealthGetHandler)

	user.Get("/", s.userGetHandler)
	user.Put("/", s.ContentTypeMiddleware(fiber.MIMEApplicationJSON), s.userPutHandler)
	user.Delete("/", s.ContentTypeMiddleware(fiber.MIMEApplicationJSON), s.userDeleteHandler)

	client := s.VersionRouter().Group("/client")
	_ = client
}

func (s *UserServer) userHealthGetHandler(c *fiber.Ctx) (err error) {
	return c.SendStatus(fiber.StatusOK)
}

func (s *UserServer) userGetHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(userServerCaller+".userGetHandler", err) }()
	parsedId, err := uuid.Parse(c.Query("id"))
	if err != nil {
		return utils.BadRequestResponse.Clone().
			WithReason("id", c.Query("id"))
	}
	if err := s.CanPass(c, parsedId); err != nil {
		return err
	}

	userType := c.Query("userType")
	if userType == "" || !ports.UserTypeValidator(userType) {
		return utils.BadRequestResponse.Clone().
			WithReason("userType", userType)
	}
	req := ports.UserUserGetRequest{
		Id:       parsedId,
		UserType: ports.UserType(userType),
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	resp, err := s.user.UserGet(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)

}

func (s *UserServer) userPutHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(userServerCaller+".authUserPostHandler", err) }()
	req := ports.UserUserPutRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse.Clone()
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	if err := s.CanPass(c, req.Id); err != nil {
		return err
	}
	resp, err := s.user.UserPut(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *UserServer) userDeleteHandler(c *fiber.Ctx) (err error) {
	defer func() { err = utils.FuncPipe(userServerCaller+".userDeleteHandler", err) }()
	req := ports.UserUserDeleteRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse.Clone()
	}
	if err := ports.Validate(c.Context(), s.Logger(), req); err != nil {
		return err
	}
	if err := s.CanPass(c, req.Id); err != nil {
		return err
	}
	userType := c.Locals("userType").(ports.UserType)
	if userType != ports.AdminUserType {
		if !ports.OtpTokenValidator(req.Token) {
			return utils.BadRequestResponse.Clone().
				WithReason("token", req.Token)
		}
		return s.user.UserDelete(c.Context(), &req, true)
	}
	return s.user.UserDelete(c.Context(), &req, false)
}
