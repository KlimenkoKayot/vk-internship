package app

import (
	"context"
	"net"

	"github.com/klimenkokayot/vk-internship/libs/logger"
	"github.com/klimenkokayot/vk-internship/service/internal/infrastructure/service"
	"github.com/klimenkokayot/vk-internship/service/internal/interfaces/grpc"
	"github.com/klimenkokayot/vk-internship/subpub/pkg/subpub"
)

type Server struct {
	grpcServer *grpc.Server
	logger     logger.Logger
}

func NewServer(logger logger.Logger) *Server {
	pubsub := subpub.NewSubPub()
	pubSubService := service.NewService(pubsub)

	return &Server{
		grpcServer: grpc.NewServer(pubSubService),
		logger:     logger,
	}
}

func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.logger.Info("Starting gRPC server", logger.Field{
		Key:   "address",
		Value: addr,
	})

	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.Stop()
	s.logger.Info("gRPC server stopped")
	return nil
}
