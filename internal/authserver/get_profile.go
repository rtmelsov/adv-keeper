package authserver

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetProfile(ctx context.Context, _ *emptypb.Empty) (*commonv1.GetProfileResponse, error) {
	log.Info("try to get profile")
	userIDRaw := ctx.Value("UserID")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		log.Error("missing user id")
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Error("error while parse the id: ", "error", err.Error())
		return nil, status.Error(codes.Unauthenticated, "can't get from string uuid")
	}
	resp, err := s.Q.GetUserByID(ctx, uid)
	if err != nil {
		log.Error("error while get user id by id: ")
		return nil, status.Error(codes.DataLoss, "нет данных по этому id")
	}

	return &commonv1.GetProfileResponse{
		Email: resp.Email,
	}, nil
}
