package api

import (
	"context"
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

func NewGRPCServer(DB *db.Service, logger *zap.Logger) *grpc.Server {
	srv := grpcServer{
		DB:     DB,
		logger: logger,
	}

	gsrv := grpc.NewServer()
	api.RegisterReceiverServiceServer(gsrv, &srv)

	reflection.Register(gsrv)

	return gsrv
}

func (s *grpcServer) CreateReceiver(ctx context.Context, req *api.CreateReceiverRequest) (*api.CreateReceiverResponse, error) {
	s.logger.Debug("call CreateReceiver method with data", zap.Any("data", req))

	if req.Receiver == nil {
		s.logger.Error("receiver has no data")
		return nil, status.Errorf(codes.InvalidArgument, "request body is empty")
	}

	receiver := &gapi.Receiver{
		Id:  uuid.New().String(),
		Url: req.Receiver.Url,
	}

	rid, err := s.DB.CreateReceiver(ctx, receiver)
	if err != nil {
		s.logger.Error("failed to create receiver", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "filed to create receiver, got error: %v", err.Error())
	}

	s.logger.Debug("receiver created successful", zap.Any("receiver_id", rid))
	return &api.CreateReceiverResponse{Rid: rid.String()}, nil
}

func (s *grpcServer) RetrieveReceiver(ctx context.Context, req *api.RetrieveReceiverRequest) (*api.RetrieveReceiverResponse, error) {
	s.logger.Debug("call RetrieveReceiver method with data", zap.Any("data", req))

	rid, err := uuid.Parse(req.Rid)
	if err != nil {
		s.logger.Error("failed to parse Rid to uuid format", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse Rid \"%v\" to uuid format, got error: %v", req.Rid, err.Error())
	}

	var receiver *gapi.Receiver
	receiver, err = s.DB.RetrieveReceiver(ctx, rid)
	if err != nil {
		s.logger.Error("failed to retrieve receiver", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "filed to retrieve receiver, got error: %v", err.Error())
	}

	s.logger.Debug("receiver retrieve successful", zap.Any("receiver", receiver))
	return &api.RetrieveReceiverResponse{Receiver: receiver}, nil
}

func (s *grpcServer) UpdateReceiver(ctx context.Context, req *api.UpdateReceiverRequest) (*api.UpdateReceiverResponse, error) {
	s.logger.Debug("call UpdateReceiver method with data", zap.Any("data", req))

	if req.Receiver == nil || req.Receiver.Url == "" || req.Receiver.Id == "" {
		s.logger.Error("invalid data to update receiver")
		return nil, status.Errorf(codes.InvalidArgument, "invalid receiver data")
	}

	err := s.DB.UpdateReceiver(ctx, req.Receiver)
	if err != nil {
		s.logger.Error("failed to update receiver", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update receiver, got error: %v", err.Error())
	}

	s.logger.Debug("receiver updated successful")
	return &api.UpdateReceiverResponse{}, nil
}

func (s *grpcServer) DeleteReceiver(ctx context.Context, req *api.DeleteReceiverRequest) (*api.DeleteReceiverResponse, error) {
	s.logger.Debug("call DeleteReceiver method with data", zap.Any("data", req))

	rid, err := uuid.Parse(req.Rid)
	if err != nil {
		s.logger.Error("failed to parse Rid to uuid format", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse Rid \"%v\" to uuid format, got error: %v", req.Rid, err.Error())
	}

	err = s.DB.DeleteReceiver(ctx, rid)
	if err != nil {
		s.logger.Error("failed to delete receiver", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete receiver, got error: %v", err.Error())
	}

	s.logger.Debug("receiver deleted successful")
	return &api.DeleteReceiverResponse{}, nil
}
