// Package middleware
package middleware

import (
	"context"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

const (
	headerAuthorize = "AccessToken"
)

var allowUnauth = map[string]struct{}{
	"/common.v1.AuthService/Register": {},
	"/common.v1.AuthService/Login":    {},
}

func ServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if _, ok := allowUnauth[info.FullMethod]; ok {
		return handler(ctx, req) // пропускаем без проверки
	}

	vals := metadata.ValueFromIncomingContext(ctx, headerAuthorize)
	if len(vals) == 0 {
		return "", status.Error(codes.Unauthenticated, "Request unauthenticated with ")
	}

	_, token, found := strings.Cut(vals[0], " ")
	if !found {
		return "", status.Error(codes.Unauthenticated, "Bad authorization string")
	}

	m, err := helpers.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	newCtx := context.WithValue(ctx, "UserID", m.UserID)

	resp, err := handler(newCtx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
