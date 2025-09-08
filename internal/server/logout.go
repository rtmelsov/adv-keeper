// Package server
package server

import (
	"context"
	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
)

func (s *Service) Logout(ctx context.Context, req *commonv1.TokenPair) (*commonv1.TokenPair, error) {
	return &commonv1.TokenPair{}, nil
}
