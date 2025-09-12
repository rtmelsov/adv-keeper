// Package authserver
package authserver

import (
	"context"
	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) Logout(ctx context.Context, _ *emptypb.Empty) (*commonv1.TokenPair, error) {
	return &commonv1.TokenPair{}, nil
}
