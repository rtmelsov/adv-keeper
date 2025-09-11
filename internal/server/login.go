// Package server
package server

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	"github.com/rtmelsov/adv-keeper/internal/helpers"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Login(ctx context.Context, req *commonv1.LoginRequest) (*commonv1.LoginResponse, error) {
	// 1) найти пользователя по email
	u, err := s.Q.GetUserByEmail(ctx, req.Email) // напиши этот sqlc-метод
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// 2) сверить пароль (bcrypt.CompareHashAndPassword)
	if bcrypt.CompareHashAndPassword([]byte(u.PwdPhc), []byte(req.Password)) != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// 3) сгенерить токены
	access, exp, err := helpers.NewAccessJWT(u.ID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "jwt")
	}

	expiresAt := timestamppb.New(exp)

	return &commonv1.LoginResponse{
		UserId: u.ID.String(),
		Email:  req.Email,
		Tokens: &commonv1.TokenPair{
			AccessToken: access,
			ExpiresAt:   expiresAt,
		},
	}, nil
}
