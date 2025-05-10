package handler

import (
	"context"

	"github.com/klimenkokayot/vk-internship/service/internal/domain"
	"github.com/klimenkokayot/vk-internship/service/pkg/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PubSubHandler struct {
	pb.UnimplementedPubSubServer
	service domain.PubSubService
}

func NewPubSubHandler(service domain.PubSubService) *PubSubHandler {
	return &PubSubHandler{
		service: service,
	}
}

func (h *PubSubHandler) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	ctx := stream.Context()
	key := req.GetKey()

	sub, err := h.service.Subscribe(ctx, key)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-sub.Chan():
			if !ok {
				return status.Error(codes.Unavailable, "subscription closed")
			}
			if err := stream.Send(&pb.Event{Data: msg}); err != nil {
				return status.Errorf(codes.Internal, "failed to send event: %v", err)
			}
		}
	}
}

func (h *PubSubHandler) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.Empty, error) {
	key := req.GetKey()
	data := req.GetData()

	if err := h.service.Publish(ctx, key, data); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish: %v", err)
	}

	return &pb.Empty{}, nil
}
