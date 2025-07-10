package user

import (
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/server"
	"github.com/kasragay/backend/internal/services"
)

const userServerCaller = packageCaller + ".UserServer"

type UserServer struct {
	*server.AbstractServer
	user ports.UserService
}

func New() ports.Server {
	s := server.NewAbstractServer(ports.UserServiceName)
	return &UserServer{
		AbstractServer: s,
		user: services.NewUserService(
			s.Logger(),
			s.Cache(),
			s.Relational(),
			s.Mongo(),
			s.S3(),
		),
	}
}
