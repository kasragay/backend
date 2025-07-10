package gateway

import (
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/server"
	"github.com/kasragay/backend/internal/services"
)

const gatewayServerCaller = packageCaller + ".GatewayServer"

type GatewayServer struct {
	*server.AbstractServer
	auth        ports.AuthService
	ratelimiter ports.RatelimiterService
}

func New() ports.Server {
	s := server.NewAbstractServer(ports.GatewayServiceName)
	return &GatewayServer{
		AbstractServer: s,
		auth: services.NewAuthService(
			s.Logger(),
			s.Cache(),
			s.Relational(),
			s.S3(),
			s.Mongo(),
			services.NewTelecomService(s.Logger()),
			services.NewMailcomService(s.Logger()),
		),
		ratelimiter: services.NewRatelimiterService(s.Logger()),
	}
}
