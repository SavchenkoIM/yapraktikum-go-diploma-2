package grpc_server

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	proto "passwordvault/internal/proto/gen"
	"passwordvault/internal/storage/server_store"
)

func (s *GRPCServer) UserLogin(ctx context.Context, in *proto.UserRequest) (*proto.UserResponse, error) {
	res := &proto.UserResponse{AccessToken: ""}

	token, err := s.dataStorage.UserLogin(ctx, in.Username, in.Password)
	if err != nil {
		msg := fmt.Sprintf("Error while logging user in: %s", err.Error())
		s.logger.Sugar().Errorf(msg)
		return res, status.Error(codes.Unauthenticated, msg)
	}
	res.AccessToken = token

	return res, nil
}

func (s *GRPCServer) UserCreate(ctx context.Context, in *proto.UserRequest) (*proto.UserResponse, error) {
	res := &proto.UserResponse{AccessToken: ""}

	err := s.dataStorage.UserRegister(ctx, in.Username, in.Password)
	if err != nil {
		if errors.Is(err, server_store.ErrUserAlreadyExists) {
			msg := fmt.Sprintf("Error while creating user: %s", err.Error())
			s.logger.Sugar().Errorf(msg)
			return res, status.Error(codes.AlreadyExists, msg)
		}

		msg := fmt.Sprintf("Error while creating user: %s", err.Error())
		s.logger.Sugar().Errorf(msg)
		return res, status.Error(codes.Unknown, msg)
	}

	token, err := s.dataStorage.UserLogin(ctx, in.Username, in.Password)
	if err != nil {
		msg := fmt.Sprintf("Error while logging user in: %s", err.Error())
		s.logger.Sugar().Errorf(msg)
		return res, status.Error(codes.Unauthenticated, msg)
	}
	res.AccessToken = token

	return res, nil
}
