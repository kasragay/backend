package server

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/bytedance/sonic"
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/repository"
	"github.com/kasragay/backend/internal/utils"
)

const serverCaller = packageCaller + ".AbstractServer"

type AbstractServer struct {
	logger      *utils.Logger
	app         *fiber.App
	verRouter   fiber.Router
	cache       ports.CacheRepo
	rel         ports.RelationalRepo
	s3          ports.S3Repo
	mongo       ports.MongoRepo
	domain      string
	version     ports.Version
	fUrl        string
	bUrl        string
	serviceName ports.ServiceName
}

func NewAbstractServer(serviceName ports.ServiceName) *AbstractServer {
	logger := utils.NewLogger(
		zap.String("service", string(serviceName)),
	)
	version_ := os.Getenv("VERSION")
	if version_ == "" {
		logger.Fatal(context.Background(), "VERSION is not set")
	}
	varsion := ports.Version(version_)
	if !slices.Contains(ports.Versions, varsion) {
		logger.Fatalf(context.Background(), "VERSION is not valid; should be one of %v", ports.Versions)
	}
	longVersion := os.Getenv("LONG_VERSION")
	if longVersion == "" {
		logger.Fatal(context.Background(), "LONG_VERSION is not set")
	}
	if !strings.HasPrefix(longVersion, version_) {
		logger.Fatalf(context.Background(), "LONG_VERSION is not valid; should start with %s", version_)
	}
	logger = utils.NewLogger(
		zap.String("service", string(serviceName)),
		zap.String("version", longVersion),
	)
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		logger.Fatal(context.Background(), "DOMAIN is not set")
	}
	logger = utils.NewLogger(
		zap.String("service", string(serviceName)),
		zap.String("version", longVersion),
		zap.String("domain", domain),
	)

	cache := repository.NewCacheRepo(logger)
	rel := repository.NewRelationalRepo(logger)
	if serviceName == ports.GatewayServiceName {
		if err := rel.TODO_DELETE_DropTables(); err != nil {
			logger.Fatalf(context.Background(), "Failed to drop tables: %v", err)
		}
		if err := rel.AutoMigrate(); err != nil {
			logger.Fatalf(context.Background(), "Failed to migrate database: %v", err)
		}
		if err := rel.TODO_DELETE_AddFirstUsers(); err != nil {
			logger.Fatalf(context.Background(), "Failed to add first users: %v", err)
		}
	}
	mongo := repository.NewMongoRepo(logger)
	s3 := repository.NewS3Repo(logger)

	fUrl := fmt.Sprintf("https://%s", domain)
	bUrl := fmt.Sprintf("https://api.%s", domain)

	s := &AbstractServer{
		app: fiber.New(fiber.Config{
			ServerHeader: fmt.Sprintf("%s - %s:%s", "Kasragay", serviceName, longVersion),
			AppName:      fmt.Sprintf("%s-%s:%s", serviceName, domain, longVersion),
			ErrorHandler: utils.ErrorHandlerFunc(logger),
			JSONEncoder:  sonic.Marshal,
			JSONDecoder:  sonic.Unmarshal,
		}),
		logger:      logger,
		cache:       cache,
		rel:         rel,
		s3:          s3,
		mongo:       mongo,
		domain:      domain,
		version:     varsion,
		fUrl:        fUrl,
		bUrl:        bUrl,
		serviceName: serviceName,
	}
	s.registerBasicRoutes()
	return s
}

func (s *AbstractServer) Logger() *utils.Logger {
	return s.logger
}

func (s *AbstractServer) App() *fiber.App {
	return s.app
}

func (s *AbstractServer) VersionRouter() fiber.Router {
	return s.verRouter
}

func (s *AbstractServer) Cache() ports.CacheRepo {
	return s.cache
}

func (s *AbstractServer) Relational() ports.RelationalRepo {
	return s.rel
}

func (s *AbstractServer) S3() ports.S3Repo {
	return s.s3
}

func (s *AbstractServer) Mongo() ports.MongoRepo {
	return s.mongo
}

func (s *AbstractServer) Domain() string {
	return s.domain
}

func (s *AbstractServer) Version() ports.Version {
	return s.version
}

func (s *AbstractServer) FrontUrl() string {
	return s.fUrl
}

func (s *AbstractServer) BackUrl() string {
	return s.bUrl
}

func (s *AbstractServer) CanPass(c *fiber.Ctx, targetId uuid.UUID, targetPhone ...string) (err error) {
	defer func() { err = utils.FuncPipe(serverCaller+".CanPass", err) }()
	sourceId := c.Locals("id").(uuid.UUID)
	sourcePhone := c.Locals("phone").(string)
	sourceUserType := c.Locals("userType").(ports.UserType)
	var ok bool
	switch {
	case sourceUserType == ports.AdminUserType:
		ok = true
	case targetPhone == nil:
		ok = sourceId == targetId
	default:
		ok = sourceId == targetId || sourcePhone == targetPhone[0]
	}
	if ok {
		return nil
	}
	return utils.JwtUnauthorizedResponse.Clone()
}
