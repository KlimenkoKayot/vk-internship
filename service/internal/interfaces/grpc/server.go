package grpc

import (
	"net"

	"github.com/klimenkokayot/vk-internship/service/internal/domain"
	"github.com/klimenkokayot/vk-internship/service/internal/interfaces/grpc/handler"
	"github.com/klimenkokayot/vk-internship/service/pkg/grpc/pb"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
}

func NewServer(pubSubService domain.PubSubService) *Server {
	srv := grpc.NewServer()
	pb.RegisterPubSubServer(srv, handler.NewPubSubHandler(pubSubService))
	return &Server{
		server: srv,
	}
}

func (s *Server) Serve(lis net.Listener) error {
	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
