package grpc_server

import (
	"context"
	proto "passwordvault/internal/proto/gen"
)

func (s *GRPCServer) DataRead(ctx context.Context, in *proto.DataReadRequest) (*proto.DataReadResponse, error) {
	res, err := s.dataStorage.DataRead(ctx, in)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *GRPCServer) DataWrite(ctx context.Context, in *proto.DataWriteRequest) (*proto.EmptyResponse, error) {
	res := &proto.EmptyResponse{}

	err := s.dataStorage.DataWrite(ctx, in)
	if err != nil {
		return nil, err
	}

	return res, nil
}
