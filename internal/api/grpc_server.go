package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	api "np_consumer/internal/api/proto"
	"np_consumer/internal/db"
)

type grpcServer struct {
	ReceiverRepository db.Service
	api.UnimplementedReceiverServiceServer
}

func NewGRPCServer(repository db.Service) *grpc.Server {
	srv := grpcServer{
		ReceiverRepository: repository,
	}

	gsrv := grpc.NewServer()
	api.RegisterReceiverServiceServer(gsrv, &srv)

	reflection.Register(gsrv)

	return gsrv
}

func (s *grpcServer) CreateReceiver(ctx context.Context, req *api.CreateReceiverRequest) (*api.CreateReceiverResponse, error) {
	fmt.Println("create receiver")
	if req.Receiver == nil {
		err := status.Errorf(codes.NotFound, "request body is empty")

		return nil, err
	}

	receiver := &db.Receiver{
		ReceiverID: uuid.New(),
		Url:        req.Receiver.Url,
	}

	rid, err := s.ReceiverRepository.CreateReceiver(ctx, receiver)
	if err != nil {
		return nil, err
	}

	return &api.CreateReceiverResponse{Rid: rid.String()}, nil
}
