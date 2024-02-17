package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	api "np_consumer/internal/api/proto"
	gapi "np_consumer/internal/api/proto"
	"np_consumer/internal/db"
)

type grpcServer struct {
	api.UnimplementedReceiverServiceServer
	DB     *db.Service
	logger *zap.Logger
}

func NewGRPCServer(DB *db.Service) *grpc.Server {
	srv := grpcServer{
		DB: DB,
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

	receiver := &gapi.Receiver{
		Id:  uuid.New().String(),
		Url: req.Receiver.Url,
	}

	rid, err := s.DB.CreateReceiver(ctx, receiver)
	if err != nil {
		return nil, err
	}

	return &api.CreateReceiverResponse{Rid: rid.String()}, nil
}

func (s *grpcServer) RetrieveReceiver(ctx context.Context, req *api.RetrieveReceiverRequest) (*api.RetrieveReceiverResponse, error) {
	fmt.Println("retrieve receiver")

	// TODO check empty data

	var receiver *gapi.Receiver
	receiver, err := s.DB.RetrieveReceiver(ctx, uuid.MustParse(req.Rid))
	if err != nil {
		return nil, err
	}

	return &api.RetrieveReceiverResponse{Receiver: receiver}, nil
}

func (s *grpcServer) DeleteReceiver(ctx context.Context, req *api.DeleteReceiverRequest) (*api.DeleteReceiverResponse, error) {
	fmt.Println("delete receiver")

	if req.Rid == "" {
		return nil, status.Errorf(codes.InvalidArgument, "the id filed if empty")
	}

	rid, err := uuid.Parse(req.Rid)
	if err != nil {
		return nil, err
	}

	err = s.DB.DeleteReceiver(ctx, rid)
	if err != nil {
		return nil, err
	}

	return &api.DeleteReceiverResponse{}, nil
}
