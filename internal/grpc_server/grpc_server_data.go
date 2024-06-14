package grpc_server

import (
	"context"
	proto "passwordvault/internal/proto/gen"
)

// Handler for Read Data request
func (s *GRPCServer) DataRead(ctx context.Context, in *proto.DataReadRequest) (*proto.DataReadResponse, error) {
	res, err := s.dataStorage.DataRead(ctx, in)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Handler for Write Data request
func (s *GRPCServer) DataWrite(ctx context.Context, in *proto.DataWriteRequest) (*proto.EmptyResponse, error) {
	res := &proto.EmptyResponse{}

	err := s.dataStorage.DataWrite(ctx, in)
	if err != nil {
		return nil, err
	}

	return res, nil
}
