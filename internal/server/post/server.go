package post

import (
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/server"
	"github.com/kasragay/backend/internal/services"
)

const postServerCaller = packageCaller + ".PostServer"

type PostServer struct {
	*server.AbstractServer
	post ports.PostService
}

func New() ports.Server {
	s := server.NewAbstractServer(ports.PostServiceName)
	return &PostServer{
		AbstractServer: s,
		post: services.NewPostService(
			s.Logger(),
			s.Cache(),
			s.Relational(),
			s.Mongo(),
			s.S3(),
		),
	}
}
