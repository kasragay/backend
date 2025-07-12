package ports

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/utils"
)

type ServiceName string

const (
	GatewayServiceName ServiceName = "gateway"
	UserServiceName    ServiceName = "user"
	PostServiceName    ServiceName = "post"
)

type Version string

const (
	V1Version     Version = "v1"
	LatestVersion Version = V1Version
)

var Versions = []Version{V1Version}

type Server interface {
	Logger() *utils.Logger
	LoggerMiddleware() fiber.Handler
	ContentTypeMiddleware(acceptedContentTypes ...string) fiber.Handler
	App() *fiber.App
	VersionRouter() fiber.Router
	Cache() CacheRepo
	Relational() RelationalRepo
	S3() S3Repo
	Mongo() MongoRepo
	Domain() string
	Version() Version
	FrontUrl() string
	BackUrl() string
	CanPass(c *fiber.Ctx, targetId uuid.UUID, targetPhone ...string) (err error)

	// Developer have to implement
	RegisterRoutes()
	RegisterV1()
}
