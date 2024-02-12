package api

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
