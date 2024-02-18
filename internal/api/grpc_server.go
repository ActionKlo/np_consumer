package api

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"np_consumer/internal/db"
	gapi "np_consumer/proto"
)

type Config struct {
}

type grpcServer struct {
	gapi.UnimplementedReceiverServiceServer
	DB     *db.PostgresService
	logger *zap.Logger
}

func NewGRPCServer(DB *db.PostgresService, logger *zap.Logger) *grpc.Server {
	srv := grpcServer{
		DB:     DB,
		logger: logger,
	}

	gsrv := grpc.NewServer()
	gapi.RegisterReceiverServiceServer(gsrv, &srv)

	reflection.Register(gsrv)

	return gsrv
}
